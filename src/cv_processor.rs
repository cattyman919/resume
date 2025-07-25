use std::{collections::HashSet, fmt::Debug, sync::Arc};

use tokio::sync::Mutex;

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

async fn get_types_per_section<T>(section: T, total_types: Arc<Mutex<HashSet<String>>>)
where
    T: IntoIterator,
    <T as IntoIterator>::Item: HasCvTypes + Debug,
{
    let mut local_types = HashSet::new();
    for item in section {
        for cv_type in item.cv_types() {
            local_types.insert(cv_type.into());
        }
    }
    total_types.lock().await.extend(local_types);
}
