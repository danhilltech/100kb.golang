use prost::Message;
use rust_bert::pipelines::keywords_extraction::{KeywordExtractionConfig, KeywordExtractionModel};
use rust_bert::pipelines::sentence_embeddings::{
    SentenceEmbeddingsBuilder, SentenceEmbeddingsConfig, SentenceEmbeddingsModelType,
};
use std::io::Cursor;
use std::ptr;

pub mod ai {
    include!(concat!(env!("OUT_DIR"), "/ai.keywords.rs"));
    include!(concat!(env!("OUT_DIR"), "/ai.sentence_embedding.rs"));
}

#[repr(C)]
pub struct SharedSentenceEmbeddingModel {
    model: rust_bert::pipelines::sentence_embeddings::SentenceEmbeddingsModel,
}

#[repr(C)]
pub struct SharedKeywordExtractionModel {
    model: rust_bert::pipelines::keywords_extraction::KeywordExtractionModel<'static>,
}

#[no_mangle]
pub extern "C" fn new_sentence_embedding() -> *mut SharedSentenceEmbeddingModel {
    let model: rust_bert::pipelines::sentence_embeddings::SentenceEmbeddingsModel =
        SentenceEmbeddingsBuilder::remote(SentenceEmbeddingsModelType::AllMiniLmL6V2)
            .with_device(tch::Device::cuda_if_available())
            .create_model()
            .unwrap();

    let m2: SharedSentenceEmbeddingModel = SharedSentenceEmbeddingModel { model: model };

    Box::into_raw(Box::new(m2))
}

#[no_mangle]
pub unsafe extern "C" fn drop_sentence_embedding(ptr: *mut SharedSentenceEmbeddingModel) {
    if ptr.is_null() {
        return;
    }
    unsafe {
        let _ = Box::from_raw(ptr);
    }
}

#[no_mangle]
pub extern "C" fn sentence_embedding(
    ptr: *mut SharedSentenceEmbeddingModel,
    req: *const libc::c_uchar,
    req_size: *mut libc::size_t,
    out_size: *mut libc::size_t,
) -> *mut u8 {
    let model = unsafe {
        assert!(!ptr.is_null());
        &mut *ptr
    };
    let bytes_raw = unsafe { std::slice::from_raw_parts(req, *req_size) };
    let bytes: Vec<u8> = Vec::from(bytes_raw);

    let sentences = ai::SentenceEmbeddingRequest::decode(&mut Cursor::new(bytes)).unwrap();

    let embd_groups: Option<Vec<Vec<f32>>> = match model.model.encode(&sentences.texts) {
        Ok(r) => Some(r),
        Err(error) => {
            println!("{}", error);
            None
        }
    };

    let mut output = ai::SentenceEmbeddingResponse::default();

    if embd_groups.is_some() {
        for group in embd_groups.unwrap().iter() {
            let mut kg = ai::Embedding::default();
            kg.vectors = group.to_vec();
            output.texts.push(kg);
        }
    }

    let mut output_vec = vec![];
    output_vec.reserve(output.encoded_len());

    output.encode(&mut output_vec).unwrap();

    output_vec.shrink_to_fit();

    let src_len = output_vec.len();

    unsafe {
        ptr::write(out_size, src_len as libc::size_t);
    }

    let slc = output_vec.into_boxed_slice();

    let res = Box::into_raw(slc);

    res.cast()
}

#[no_mangle]
pub extern "C" fn new_keyword_extraction() -> *mut SharedKeywordExtractionModel {
    let keyword_extraction_config = KeywordExtractionConfig {
        sentence_embeddings_config: SentenceEmbeddingsConfig::from(
            SentenceEmbeddingsModelType::AllMiniLmL6V2,
        ),
        // scorer_type: KeywordScorerType::MaxSum,
        ngram_range: (1, 1),
        num_keywords: 3,
        ..Default::default()
    };

    let keyword_extraction_model = KeywordExtractionModel::new(keyword_extraction_config).unwrap();

    let m2 = SharedKeywordExtractionModel {
        model: keyword_extraction_model,
    };

    Box::into_raw(Box::new(m2))
}

#[no_mangle]
pub unsafe extern "C" fn drop_keyword_extraction(ptr: *mut SharedKeywordExtractionModel) {
    if ptr.is_null() {
        return;
    }
    unsafe {
        let _ = Box::from_raw(ptr);
    }
}

#[no_mangle]
pub extern "C" fn keyword_extraction(
    ptr: *mut SharedKeywordExtractionModel,
    req: *const libc::c_uchar,
    req_size: *mut libc::size_t,
    out_size: *mut libc::size_t,
) -> *mut u8 {
    let model = unsafe {
        assert!(!ptr.is_null());
        &mut *ptr
    };

    let bytes_raw = unsafe { std::slice::from_raw_parts(req, *req_size) };
    let bytes: Vec<u8> = Vec::from(bytes_raw);

    let sentences = ai::KeywordRequest::decode(&mut Cursor::new(bytes)).unwrap();

    let key_groups: Option<Vec<Vec<rust_bert::pipelines::keywords_extraction::Keyword>>> =
        match model.model.predict(&sentences.texts) {
            Ok(r) => Some(r),
            Err(error) => {
                println!("{}", error);
                None
            }
        };

    let mut output = ai::KeywordResponse::default();

    if key_groups.is_some() {
        for group in key_groups.unwrap().iter() {
            let mut kg = ai::Keywords::default();

            for keyword in group.iter() {
                let mut kw = ai::Keyword::default();
                kw.text = keyword.text.clone().into_bytes();
                kw.score = keyword.score.clone();
                kg.keywords.push(kw);
            }
            output.texts.push(kg);
        }
    }

    let mut output_vec = vec![];

    output_vec.reserve(output.encoded_len());

    output.encode(&mut output_vec).unwrap();

    output_vec.shrink_to_fit();

    let src_len = output_vec.len();

    unsafe {
        ptr::write(out_size, src_len as libc::size_t);
    }

    let slc = output_vec.into_boxed_slice();

    let res = Box::into_raw(slc);

    res.cast()
}

#[no_mangle]
pub unsafe extern "C" fn drop_bytesarray(ptr: *mut u8) {
    if ptr.is_null() {
        return;
    }
    unsafe {
        let _ = Box::from_raw(ptr);
    }
}
