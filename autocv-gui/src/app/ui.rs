use autocv_core::cv_model::PersonalInfo;

use crate::{
    actor::State,
    app::{App, AppTab},
};

pub fn general_ui(ui: &mut egui::Ui, local_state: &mut State) {
    personal_info_ui(ui, &mut local_state.general_cv.personal_info);
    skills_achivements_ui(ui);
    education_ui(ui);
    awards_ui(ui);
}

pub fn personal_info_ui(ui: &mut egui::Ui, personal_info: &mut PersonalInfo) {
    ui.heading("Personal Info");
    ui.label("Name");
    ui.add(egui::TextEdit::singleline(&mut personal_info.name).hint_text("Your Name"));
    ui.label("Location");
    ui.add(egui::TextEdit::singleline(&mut personal_info.location).hint_text("Your Location"));
    ui.label("Email");
    ui.add(egui::TextEdit::singleline(&mut personal_info.email).hint_text("Your email"));
    ui.label("Phone");
    ui.add(egui::TextEdit::singleline(&mut personal_info.phone).hint_text("Your phone number"));
    ui.label("Website");
    ui.add(egui::TextEdit::singleline(&mut personal_info.website).hint_text("Your website"));
    ui.label("LinkedIn");
    ui.add(
        egui::TextEdit::singleline(&mut personal_info.linkedin).hint_text("Your LinkedIn profile"),
    );
    ui.label("LinkedIn Handle");
    ui.add(
        egui::TextEdit::singleline(&mut personal_info.linkedin_handle)
            .hint_text("Your LinkedIn handle"),
    );
    ui.label("Github");
    ui.add(egui::TextEdit::singleline(&mut personal_info.github).hint_text("Your GitHub profile"));
    ui.label("Github Handle");
    ui.add(
        egui::TextEdit::singleline(&mut personal_info.github_handle)
            .hint_text("Your GitHub handle"),
    );
    ui.label("Profile Pic");
    ui.add(egui::TextEdit::singleline(&mut personal_info.profile_pic).hint_text("Profile Pic"));
}

pub fn skills_achivements_ui(ui: &mut egui::Ui) {
    ui.heading("Skills & Achievements");
    ui.label("Skills");
    ui.add(egui::TextEdit::multiline(&mut String::new()).hint_text("Your Skills"));
    ui.label("Achievements");
    ui.add(egui::TextEdit::multiline(&mut String::new()).hint_text("Your Achievements"));
}

pub fn education_ui(ui: &mut egui::Ui) {
    ui.heading("Education");
    ui.label("Institution");
    ui.add(egui::TextEdit::singleline(&mut String::new()).hint_text("Your Institution"));
    ui.label("Degree");
    ui.add(egui::TextEdit::singleline(&mut String::new()).hint_text("Your Degree"));
    ui.label("Dates");
    ui.add(egui::TextEdit::singleline(&mut String::new()).hint_text("Your Dates"));
    ui.label("GPA");
    ui.add(egui::TextEdit::singleline(&mut String::new()).hint_text("Your GPA"));
    ui.label("Details");
    ui.add(egui::TextEdit::singleline(&mut String::new()).hint_text("Your GPA"));
}

pub fn awards_ui(ui: &mut egui::Ui) {
    ui.heading("Awards");
    ui.label("Title");
    ui.add(egui::TextEdit::singleline(&mut String::new()).hint_text("Title"));
    ui.label("Organization");
    ui.add(egui::TextEdit::singleline(&mut String::new()).hint_text("Organization"));
    ui.label("Date");
    ui.add(egui::TextEdit::singleline(&mut String::new()).hint_text("Date"));
    ui.label("Points");
    ui.add(egui::TextEdit::multiline(&mut String::new()).hint_text("Points"));
}

pub fn project_ui(ui: &mut egui::Ui) {
    ui.label("Name");
    ui.add(egui::TextEdit::singleline(&mut String::new()).hint_text("Project Name"));
    ui.label("Github");
    ui.add(egui::TextEdit::singleline(&mut String::new()).hint_text("Project Github"));
    ui.label("Github Handle");
    ui.add(egui::TextEdit::singleline(&mut String::new()).hint_text("Project Github Handle"));
    ui.label("Description");
    ui.add(egui::TextEdit::multiline(&mut String::new()).hint_text("Project Description"));
    ui.label("CV Types");
    ui.add(egui::TextEdit::multiline(&mut String::new()).hint_text("Project CV Types"));
    ui.label("Points");
    ui.add(egui::TextEdit::multiline(&mut String::new()).hint_text("Project Points"));
}

pub fn experience_ui(ui: &mut egui::Ui) {
    ui.label("Company");
    ui.add(egui::TextEdit::singleline(&mut String::new()).hint_text("Experience Company"));
    ui.label("Location");
    ui.add(egui::TextEdit::singleline(&mut String::new()).hint_text("Experience Location"));
    ui.label("Role");
    ui.add(egui::TextEdit::singleline(&mut String::new()).hint_text("Experience Role"));
    ui.label("Dates");
    ui.add(egui::TextEdit::singleline(&mut String::new()).hint_text("Experience Dates"));
    ui.label("Job Type");
    ui.add(egui::TextEdit::multiline(&mut String::new()).hint_text("Experience CV Types"));
    ui.label("CV Types");
    ui.add(egui::TextEdit::multiline(&mut String::new()).hint_text("Experience CV Types"));
    ui.label("Points");
    ui.add(egui::TextEdit::multiline(&mut String::new()).hint_text("Experience Points"));
}

pub fn side_panel_ui(ctx: &egui::Context, app: &mut App) {
    egui::SidePanel::left("side_panel")
        .min_width(150.0)
        .show(ctx, |ui| {
            ui.heading("Config Tab");
            ui.separator();
            ui.add_space(10.0);
            ui.spacing_mut().item_spacing.y = 10.0;
            // Get the style once to reuse it
            let theme = &ctx.style().visuals.clone();
            let selected_fill = theme.selection.bg_fill;
            let selected_stroke = theme.selection.stroke;

            // -- General Tab Button --
            let is_general_selected = app.selected_tab == AppTab::General;
            let general_button = egui::Button::new("General")
                .min_size(egui::vec2(120.0, 25.0)) // Give it a fixed size
                .fill(if is_general_selected {
                    selected_fill
                } else {
                    egui::Color32::TRANSPARENT
                })
                .stroke(if is_general_selected {
                    selected_stroke
                } else {
                    egui::Stroke::NONE
                });

            if ui.add(general_button).clicked() {
                app.selected_tab = AppTab::General;
            }

            // -- Experiences Tab Button --
            let is_experiences_selected = app.selected_tab == AppTab::Experiences;
            let experiences_button = egui::Button::new("Experiences")
                .min_size(egui::vec2(120.0, 25.0))
                .fill(if is_experiences_selected {
                    selected_fill
                } else {
                    egui::Color32::TRANSPARENT
                })
                .stroke(if is_experiences_selected {
                    selected_stroke
                } else {
                    egui::Stroke::NONE
                });

            if ui.add(experiences_button).clicked() {
                app.selected_tab = AppTab::Experiences;
            }

            // -- Projects Tab Button --
            let is_projects_selected = app.selected_tab == AppTab::Projects;
            let projects_button = egui::Button::new("Projects")
                .min_size(egui::vec2(120.0, 25.0))
                .fill(if is_projects_selected {
                    selected_fill
                } else {
                    egui::Color32::TRANSPARENT
                })
                .stroke(if is_projects_selected {
                    selected_stroke
                } else {
                    egui::Stroke::NONE
                });

            if ui.add(projects_button).clicked() {
                app.selected_tab = AppTab::Projects;
            }
        });
}
