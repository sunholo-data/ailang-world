use std::{fs, path::Path};

// Re-embed the frontend whenever any file under dist/ changes: a fingerprint of the files
// goes into an env var the crate reads, so cargo recompiles main.rs and generate_context!.
fn walk(dir: &Path, acc: &mut u64) {
    if let Ok(rd) = fs::read_dir(dir) {
        for e in rd.flatten() {
            let p = e.path();
            println!("cargo:rerun-if-changed={}", p.display());
            if p.is_dir() { walk(&p, acc); continue; }
            if let Ok(b) = fs::read(&p) {
                for (i, x) in b.iter().enumerate() { *acc = acc.wrapping_mul(1099511628211).wrapping_add(*x as u64 ^ i as u64); }
            }
        }
    }
}

fn main() {
    let mut h: u64 = 1469598103934665603;
    println!("cargo:rerun-if-changed=dist");
    walk(Path::new("dist"), &mut h);
    println!("cargo:rustc-env=WORLD_DIST_FINGERPRINT={h:x}");
    tauri_build::build()
}
