#![warn(clippy::all, rust_2018_idioms)]
// #![cfg_attr(not(debug_assertions), windows_subsystem = "windows")] // hide console window on Windows in release

mod actor;
mod app;

use crate::actor::State;
use app::App;
use std::sync::{Arc, Mutex};
use tokio::sync::mpsc;

// When compiling natively:
#[cfg(not(target_arch = "wasm32"))]
#[tokio::main]
async fn main() -> Result<(), eframe::Error> {
    env_logger::init();

    let (sender, receiver) = mpsc::channel(32);

    let shared_state = Arc::new(Mutex::new(State::default()));

    // Actor Initialization
    let actor_shared_state = Arc::clone(&shared_state);
    tokio::spawn(async move {
        let actor = actor::Actor::new(receiver, actor_shared_state);
        actor.run().await;
    });

    // Initial Load CV Task
    let sender_clone = sender.clone();
    tokio::spawn(async move {
        if let Err(e) = sender_clone.send(actor::ActorMessage::LoadCV()).await {
            log::error!("Failed to send message to actor: {}", e);
        }
    });

    let native_options = eframe::NativeOptions {
        viewport: egui::ViewportBuilder::default()
            .with_inner_size([400.0, 300.0])
            .with_fullscreen(false)
            .with_min_inner_size([300.0, 220.0]),
        ..Default::default()
    };

    eframe::run_native(
        "Auto CV GUI",
        native_options,
        Box::new(|cc| Ok(Box::new(crate::App::new(cc, sender, shared_state)))),
    )?;

    Ok(())
}
