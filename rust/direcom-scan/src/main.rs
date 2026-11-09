use std::{env, fs, path::{Path, PathBuf}, process};

fn walk(path: &Path, root: &Path, files: &mut u64, dirs: &mut u64, bytes: &mut u64) -> std::io::Result<()> {
    for item in fs::read_dir(path)? {
        let item = item?;
        let item_path = item.path();
        let metadata = fs::symlink_metadata(&item_path)?;
        if metadata.file_type().is_symlink() { continue; }
        if metadata.is_dir() {
            *dirs += 1;
            walk(&item_path, root, files, dirs, bytes)?;
        } else if metadata.is_file() {
            *files += 1;
            *bytes += metadata.len();
        }
        let _ = root;
    }
    Ok(())
}

fn json_escape(value: &str) -> String {
    value.chars().flat_map(|ch| match ch {
        '"' => "\\\"".chars().collect::<Vec<_>>(),
        '\\' => "\\\\".chars().collect(),
        '\n' => "\\n".chars().collect(),
        '\r' => "\\r".chars().collect(),
        '\t' => "\\t".chars().collect(),
        c if c.is_control() => format!("\\u{:04x}", c as u32).chars().collect(),
        c => vec![c],
    }).collect()
}

fn main() {
    let root = env::args_os().nth(1).map(PathBuf::from).unwrap_or_else(|| PathBuf::from("."));
    let root = match fs::canonicalize(&root) {
        Ok(path) if path.is_dir() => path,
        Ok(_) => { eprintln!("path is not a directory"); process::exit(2); }
        Err(error) => { eprintln!("cannot access path: {error}"); process::exit(2); }
    };
    let mut files = 0;
    let mut dirs = 0;
    let mut bytes = 0;
    if let Err(error) = walk(&root, &root, &mut files, &mut dirs, &mut bytes) {
        eprintln!("scan failed: {error}"); process::exit(1);
    }
    println!("{{\"root\":\"{}\",\"files\":{},\"directories\":{},\"bytes\":{}}}", json_escape(&root.display().to_string()), files, dirs, bytes);
}
