use std::io::Result;
fn main() -> Result<()> {
    prost_build::compile_protos(
        &[
            "src/keywords.proto",
            "src/sentence_embedding.proto",
            "src/zero_shot.proto",
        ],
        &["src/"],
    )?;
    Ok(())
}
