#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

// Forces a recompile (and so a fresh frontend embed) when dist/ changes; see build.rs.
const _DIST: &str = env!("WORLD_DIST_FINGERPRINT");

fn main() {
    tauri::Builder::default()
        .run(tauri::generate_context!())
        .expect("error while running World");
}
