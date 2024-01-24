use adblock::{
    lists::{FilterSet, ParseOptions},
    request::Request,
    Engine,
};
use prost::Message;
use std::collections::HashSet;
use std::io::Cursor;
use std::ptr;

pub mod goadblock {
    include!(concat!(env!("OUT_DIR"), "/adblock.content.rs"));
}

#[repr(C)]
pub struct AdblockEngine {
    engine: Engine,
}

#[no_mangle]
pub extern "C" fn new_adblock(
    req: *const libc::c_uchar,
    req_size: *mut libc::size_t,
) -> *mut AdblockEngine {
    let bytes_raw = unsafe { std::slice::from_raw_parts(req, *req_size) };
    let bytes: Vec<u8> = Vec::from(bytes_raw);

    let request = goadblock::Rules::decode(&mut Cursor::new(bytes)).unwrap();

    let debug_info = false;
    let mut filter_set = FilterSet::new(debug_info);
    filter_set.add_filters(&request.rules, ParseOptions::default());

    let engine = Engine::from_filter_set(filter_set, true);

    let m2: AdblockEngine = AdblockEngine { engine: engine };

    Box::into_raw(Box::new(m2))
}

#[no_mangle]
pub unsafe extern "C" fn drop_adblock(ptr: *mut AdblockEngine) {
    if ptr.is_null() {
        return;
    }
    unsafe {
        let _ = Box::from_raw(ptr);
    }
}

#[no_mangle]
pub extern "C" fn filter(
    ptr: *mut AdblockEngine,
    req: *const libc::c_uchar,
    req_size: *mut libc::size_t,
    out_size: *mut libc::size_t,
) -> *mut u8 {
    let engine = unsafe {
        assert!(!ptr.is_null());
        &mut *ptr
    };

    let bytes_raw = unsafe { std::slice::from_raw_parts(req, *req_size) };
    let bytes: Vec<u8> = Vec::from(bytes_raw);

    let request = goadblock::FilterRequest::decode(&mut Cursor::new(bytes)).unwrap();

    let matches = engine.engine.hidden_class_id_selectors(
        &request.classes,
        &request.ids,
        &HashSet::default(),
    );

    let mut output = goadblock::FilterResponse::default();

    for url in request.urls.iter() {
        let request: Option<Request> = match Request::new(url, &request.base_url, "other") {
            Ok(r) => Some(r),
            Err(_error) => {
                println!("URL ERROR: {:?}", url);
                None
            }
        };

        if request.is_some() {
            let blocked = engine.engine.check_network_request(&request.unwrap());
            if blocked.matched {
                output.blocked_domains.push(url.to_string());
            }
        }
    }

    output.matches = matches;

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
