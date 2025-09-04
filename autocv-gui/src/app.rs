mod style;
mod ui;

use std::sync::{Arc, Mutex};

use log::{info, warn};
use tokio::sync::mpsc;
use ui::general_ui;

use crate::{
    actor::{ActorMessage, State},
    app::{
        style::apply_style,
        ui::{experience_ui, project_ui, side_panel_ui, top_panel_ui},
    },
};

#[derive(PartialEq)]
pub enum AppTab {
    General,
    Experiences,
    Projects,
}

pub struct App {
    selected_tab: AppTab,
    shared_state: Arc<Mutex<State>>,
    local_state: State,
    actor_sender: mpsc::Sender<ActorMessage>,
}

impl App {
    /// Called once before the first frame.
    pub fn new(
        cc: &eframe::CreationContext<'_>,
        actor_sender: mpsc::Sender<ActorMessage>,
        shared_state: Arc<Mutex<State>>,
    ) -> Self {
        // This is also where you can customize the look and feel of egui using
        // `cc.egui_ctx.set_visuals` and `cc.egui_ctx.set_fonts`.
        // Default::default()

        let style = apply_style();
        cc.egui_ctx.set_style(style);

        Self {
            selected_tab: AppTab::General,
            local_state: Arc::clone(&shared_state).lock().unwrap().clone(),
            shared_state,
            actor_sender,
        }
    }
}

impl eframe::App for App {
    /// Called by the framework to save state before shutdown.
    fn save(&mut self, _storage: &mut dyn eframe::Storage) {
        // eframe::set_value(storage, eframe::APP_KEY, self);
    }

    /// Called each time the UI needs repainting, which may be many times per second.
    fn update(&mut self, ctx: &egui::Context, _frame: &mut eframe::Frame) {
        // For inspiration and more examples, go to https://emilk.github.io/egui

        // Top Panel
        let top_panel_frame = style::top_panel_frame();
        top_panel_ui(ctx, top_panel_frame);

        // Side Panel
        side_panel_ui(ctx, self);

        // Central Panel
        egui::CentralPanel::default().show(ctx, |ui| {
            ui.heading("Configuration");
            ui.separator();

            egui::ScrollArea::vertical().show(ui, |ui| {
                match self.selected_tab {
                    AppTab::General => general_ui(ui, &mut self.local_state.general_cv),
                    AppTab::Experiences => {
                        for (i, experience) in
                            self.local_state.experiences_cv.iter_mut().enumerate()
                        {
                            experience_ui(ui, experience, i);
                            ui.separator();
                        }
                    }
                    AppTab::Projects => {
                        for (i, project) in self.local_state.projects_cv.iter_mut().enumerate() {
                            project_ui(ui, project, i);
                            ui.separator();
                        }
                    }
                }
                ui.allocate_space(ui.available_size());
            });

            // if let Ok(state) = self.shared_state.try_lock() {
            //     info!("Successfully locked the shared state.",);
            //     ui.label(&state.general_cv.personal_info.name);
            // } else {
            //     warn!("Failed to lock the shared state.");
            //     ui.label("State is being updated...");
            // }
            //
            // Lock the mutex to safely read the data
            // let state = self.shared_state.lock().unwrap();
            // ui.label(&state);

            // ctx.request_repaint(); // Ensures the UI updates when the actor changes the state
        });
    }
}
