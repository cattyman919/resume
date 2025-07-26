use std::{collections::HashSet, path::Path, sync::Arc};

use tokio::{fs, io, sync::Mutex};

use crate::types::{CVData, HasCvTypes};
pub async fn get_all_cv_types(cv_data: &CVData) -> Result<HashSet<String>, &'static str> {
    let total_types = Arc::new(Mutex::new(HashSet::<String>::new()));

    let experience_task = get_types_per_section(&cv_data.experiences, Arc::clone(&total_types));
    let project_task = get_types_per_section(&cv_data.projects, Arc::clone(&total_types));

    tokio::join!(experience_task, project_task);

    let final_set = Arc::try_unwrap(total_types)
        .expect("Error: Arc should be uniquely owned here")
        .into_inner();

    Ok(final_set)
}

async fn get_types_per_section<T>(section: &[T], total_types: Arc<Mutex<HashSet<String>>>)
where
    T: HasCvTypes,
{
    let mut local_types = HashSet::new();
    for item in section {
        for cv_type in item.cv_types() {
            local_types.insert(cv_type.into());
        }
    }
    total_types.lock().await.extend(local_types);
}

pub async fn write_cv(cv_type: &str) -> io::Result<()> {
    let main_path = format!("cv/{cv_type}/main_cv");
    let bw_path = format!("cv/{cv_type}/bw_cv");

    copy_dir_recursively("template_cv/main_cv", main_path).await?;
    copy_dir_recursively("template_cv/bw_cv", bw_path).await?;

    Ok(())
}

async fn copy_dir_recursively<P: AsRef<Path>, Q: AsRef<Path>>(from: P, to: Q) -> io::Result<()> {
    let from_path = from.as_ref();
    let to_path = to.as_ref();

    fs::create_dir_all(&to_path).await?;

    let mut reader = fs::read_dir(from_path).await?;

    while let Some(entry) = reader.next_entry().await? {
        let file_type = entry.file_type().await?;

        let to_path_entry = to_path.join(entry.file_name());

        if file_type.is_dir() {
            Box::pin(copy_dir_recursively(entry.path(), &to_path_entry)).await?;
        } else {
            fs::copy(entry.path(), &to_path_entry).await?;
        }
    }

    Ok(())
}
