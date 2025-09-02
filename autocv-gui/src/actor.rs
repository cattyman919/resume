use std::{
    error::Error,
    sync::{Arc, Mutex},
};

use anyhow::Result;
use autocv_core::{
    cv_model::{ExperiencesCVData, GeneralCVData, ProjectsCVData},
    cv_processor, load_cv,
};
use log::info;
use tokio::sync::mpsc;

#[derive(Debug)]
pub enum ActorMessage {
    LoadCV(),
}

#[derive(Clone, Default)]
pub struct State {
    pub general_cv: GeneralCVData,
    pub experiences_cv: ExperiencesCVData,
    pub projects_cv: ProjectsCVData,
}

pub struct Actor {
    receiver: mpsc::Receiver<ActorMessage>,
    shared_state: Arc<Mutex<State>>,
}

impl Actor {
    pub fn new(receiver: mpsc::Receiver<ActorMessage>, shared_state: Arc<Mutex<State>>) -> Self {
        Self {
            receiver,
            shared_state,
        }
    }

    pub async fn run(mut self) {
        while let Some(message) = self.receiver.recv().await {
            let _result = self.handle_message(message).await;
        }
    }

    async fn handle_message(&mut self, msg: ActorMessage) {
        info!("Actor received message: {:?}", msg);
        match msg {
            ActorMessage::LoadCV() => {
                setup(self).await.unwrap();
            }
        }
    }
}

pub async fn setup(actor: &mut Actor) -> Result<(), Box<dyn Error>> {
    cv_processor::setup_directories().await?;

    info!("Loading YAML Data...");
    let (general_cv, projects_cv, experiences_cv) = load_cv::load_cv_data().await?;

    {
        let mut shared_state = actor.shared_state.lock().unwrap();
        shared_state.general_cv = general_cv;
        shared_state.projects_cv = projects_cv;
        shared_state.experiences_cv = experiences_cv;
    }

    // info!("Getting Total CV Types...");
    // let all_cv_types = cv_processor::get_all_cv_types(&projects_cv, &experiences_cv).await?;
    // info!("All CV Types: {:?}", all_cv_types);

    Ok(())
}

// pub async fn generate_all_cvs() -> Result<(), Box<dyn Error>> {
//     let processing_tasks = all_cv_types.into_iter().map(|cv_type| {
//         let general_cv_clone = Arc::clone(&general_cv);
//         let projects_cv_clone = projects_cv.as_ref().clone();
//         let experiences_cv_clone = experiences_cv.as_ref().clone();
//         tokio::spawn(async move {
//             println!("Processing CV type: {}", cv_type);
//             cv_processor::write_cv(
//                 general_cv_clone,
//                 projects_cv_clone,
//                 experiences_cv_clone,
//                 cv_type,
//                 false,
//             )
//             .await
//         })
//     });
//
//     let results = join_all(processing_tasks).await;
//
//     // --- Error Handling for Concurrent Tasks ---
//     let mut errors = Vec::new();
//     for result in results {
//         match result {
//             // The task itself panicked (a serious bug).
//             Err(join_error) => errors.push(format!("A task panicked: {}", join_error)),
//             // The task completed but returned an error.
//             Ok(Err(task_error)) => errors.push(format!("A task failed: {}", task_error)),
//             Ok(Ok(_)) => (),
//         }
//     }
//
//     if !errors.is_empty() {
//         eprintln!("\nErrors occurred during CV generation:");
//         for e in errors {
//             eprintln!("- {}", e);
//         }
//         return Err("CV generation failed due to one or more task errors.".into());
//     }
//
//     cv_processor::move_aux_files().await?;
//
//     println!("\n==== All LaTeX CV Generation Complete ====");
//     Ok(())
// }
