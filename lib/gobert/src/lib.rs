use rust_bert::pipelines::sentence_embeddings::{
    SentenceEmbeddingsBuilder, SentenceEmbeddingsModelType,
};
use std::{ffi::CStr, ptr};

#[repr(C)]
pub struct Model {
    model: rust_bert::pipelines::sentence_embeddings::SentenceEmbeddingsModel,
}

#[no_mangle]
pub extern "C" fn new_sentence_embedding() -> *mut Model {
    let model: rust_bert::pipelines::sentence_embeddings::SentenceEmbeddingsModel =
        SentenceEmbeddingsBuilder::remote(SentenceEmbeddingsModelType::AllMiniLmL6V2)
            .with_device(tch::Device::cuda_if_available())
            .create_model()
            .unwrap();

    let m2 = Model { model: model };

    Box::into_raw(Box::new(m2))
}

#[no_mangle]
pub unsafe extern "C" fn drop_sentence_embedding(ptr: *mut Model) {
    if ptr.is_null() {
        return;
    }
    unsafe {
        let _ = Box::from_raw(ptr);
    }
}

#[no_mangle]
pub extern "C" fn sentence_embedding(
    ptr: *mut Model,
    strs: *const *const libc::c_char,
    dst_ptr: *mut f32,
) {
    let model = unsafe {
        assert!(!ptr.is_null());
        &mut *ptr
    };
    // Define input
    let sentences = ["this is an example sentence"];

    unsafe {
        for i in 0.. {
            let member_ptr: *const libc::c_char = *(strs.offset(i));
            if member_ptr != ptr::null() {
                let member: &CStr = CStr::from_ptr(member_ptr);
                println!("d   {}", member.to_str().unwrap());
            } else {
                break;
            }
        }
    }

    let embeddings = model.model.encode(&sentences).unwrap();

    // Box::into_raw(Box::new(embeddings))

    let mut flat_embeddings: Vec<f32> = embeddings.into_iter().flatten().collect();

    flat_embeddings.shrink_to_fit();
    assert!(flat_embeddings.len() == flat_embeddings.capacity());

    let src_ptr = flat_embeddings.as_mut_ptr();
    let src_len = flat_embeddings.len();

    unsafe {
        ptr::copy_nonoverlapping(src_ptr, dst_ptr, src_len);
    }
}
