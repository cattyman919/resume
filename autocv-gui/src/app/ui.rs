use autocv_core::cv_model::{
    Award, Certificate, Education, Experience, GeneralCVConfig, PersonalInfo, Project,
    SkillsAchievements,
};
use log::info;

use crate::{
    actor::ActorMessage,
    app::{App, AppTab},
};

// A helper function to create a consistent section header
fn section_header(ui: &mut egui::Ui, title: &str) {
    ui.add_space(10.0);
    ui.heading(title);
    ui.add_space(4.0);
}

// A helper for adding/removing items from a Vec<String>
enum TextEditMode {
    SingleLine,
    MultiLine,
}

fn editable_list(
    ui: &mut egui::Ui,
    title: &str,
    text_edit_mode: TextEditMode,
    items: &mut Vec<String>,
) {
    ui.label(title);
    let mut item_to_remove = None;
    for (i, item) in items.iter_mut().enumerate() {
        ui.horizontal(|ui| {
            match text_edit_mode {
                TextEditMode::SingleLine => {
                    ui.add(egui::TextEdit::singleline(item));
                }
                TextEditMode::MultiLine => {
                    ui.add(
                        egui::TextEdit::multiline(item).desired_width(ui.available_width() * 0.9),
                    );
                }
            }
            if ui
                .add(egui::Button::new("➖").frame(false))
                .on_hover_text("Remove")
                .clicked()
            {
                item_to_remove = Some(i);
            }
        });
    }

    if let Some(index) = item_to_remove {
        items.remove(index);
    }

    // Use a secondary-styled button for adding
    if ui.add(egui::Button::new("➕ Add")).clicked() {
        items.push(String::new());
    }
    ui.add_space(5.0);
}

// A specific helper for the Vec<Certificate>
fn certificate_list(ui: &mut egui::Ui, title: &str, items: &mut Vec<Certificate>) {
    ui.label(title);
    let mut item_to_remove = None;
    for (i, cert) in items.iter_mut().enumerate() {
        ui.horizontal(|ui| {
            ui.add(egui::TextEdit::singleline(&mut cert.year).desired_width(40.0));
            ui.add(
                egui::TextEdit::singleline(&mut cert.name)
                    .desired_width(ui.available_width() * 0.9),
            );
            if ui
                .add(egui::Button::new("➖").frame(false))
                .on_hover_text("Remove")
                .clicked()
            {
                item_to_remove = Some(i);
            }
        });
    }

    if let Some(index) = item_to_remove {
        items.remove(index);
    }

    if ui.button("➕ Add Certificate").clicked() {
        items.push(Certificate::default());
    }
    ui.add_space(5.0);
}

pub fn general_ui(ui: &mut egui::Ui, general_cv: &mut GeneralCVConfig) {
    personal_info_ui(ui, &mut general_cv.personal_info);
    skills_achivements_ui(ui, &mut general_cv.skills_achievements);

    section_header(ui, "Education");
    // TODO: Add a button to add/remove education entries
    for (i, education) in general_cv.education.iter_mut().enumerate() {
        education_ui(ui, education, i);
    }

    section_header(ui, "Awards");
    // TODO: Add a button to add/remove award entries
    for (i, award) in general_cv.awards.iter_mut().enumerate() {
        awards_ui(ui, award, i);
    }
}

pub fn personal_info_ui(ui: &mut egui::Ui, personal_info: &mut PersonalInfo) {
    section_header(ui, "Personal Info");
    egui::Grid::new("personal_info_grid")
        .num_columns(2)
        .spacing([40.0, 4.0])
        .striped(true)
        .show(ui, |ui| {
            ui.label("Name");
            ui.text_edit_singleline(&mut personal_info.name);
            ui.end_row();

            ui.label("Location");
            ui.text_edit_singleline(&mut personal_info.location);
            ui.end_row();

            ui.label("Email");
            ui.text_edit_singleline(&mut personal_info.email);
            ui.end_row();

            ui.label("Phone");
            ui.text_edit_singleline(&mut personal_info.phone);
            ui.end_row();

            ui.label("Website");
            ui.text_edit_singleline(&mut personal_info.website);
            ui.end_row();

            ui.label("LinkedIn Profile");
            ui.text_edit_singleline(&mut personal_info.linkedin);
            ui.end_row();

            ui.label("GitHub Profile");
            ui.text_edit_singleline(&mut personal_info.github);
            ui.end_row();
        });
}

pub fn skills_achivements_ui(ui: &mut egui::Ui, skills_achievements: &mut SkillsAchievements) {
    section_header(ui, "Skills & Achievements");

    egui::CollapsingHeader::new("Skills & Certificates").show(ui, |ui| {
        editable_list(
            ui,
            "Hard Skills",
            TextEditMode::SingleLine,
            &mut skills_achievements.hard_skills,
        );
        editable_list(
            ui,
            "Soft Skills",
            TextEditMode::SingleLine,
            &mut skills_achievements.soft_skills,
        );
        editable_list(
            ui,
            "Programming Languages",
            TextEditMode::SingleLine,
            &mut skills_achievements.programming_languages,
        );
        editable_list(
            ui,
            "Databases",
            TextEditMode::SingleLine,
            &mut skills_achievements.databases,
        );
        editable_list(
            ui,
            "Misc Tools",
            TextEditMode::SingleLine,
            &mut skills_achievements.misc,
        );
        certificate_list(ui, "Certificates", &mut skills_achievements.certificates);
    });
}

pub fn education_ui(ui: &mut egui::Ui, education: &mut Education, i: usize) {
    egui::CollapsingHeader::new(if education.institution.is_empty() {
        "New Entry"
    } else {
        &education.institution
    })
    .show(ui, |ui| {
        egui::Grid::new(format!("education_{}", i))
            .num_columns(2)
            .spacing([40.0, 4.0])
            .striped(true)
            .show(ui, |ui| {
                ui.label("Institution");
                ui.text_edit_singleline(&mut education.institution);
                ui.end_row();

                ui.label("Degree");
                ui.text_edit_singleline(&mut education.degree);
                ui.end_row();

                ui.label("Dates");
                ui.text_edit_singleline(&mut education.dates);
                ui.end_row();

                ui.label("GPA");
                ui.text_edit_singleline(&mut education.gpa);
                ui.end_row();
            });

        ui.add_space(10.0);
        editable_list(
            ui,
            "Details",
            TextEditMode::MultiLine,
            &mut education.details,
        );
    });
}

pub fn awards_ui(ui: &mut egui::Ui, award: &mut Award, i: usize) {
    egui::CollapsingHeader::new(if award.title.is_empty() {
        "New Award".to_string()
    } else {
        award.title.clone()
    })
    .show(ui, |ui| {
        egui::Grid::new(format!("award_{}", i))
            .num_columns(2)
            .spacing([40.0, 4.0])
            .striped(true)
            .show(ui, |ui| {
                ui.label("Title");
                ui.text_edit_singleline(&mut award.title);
                ui.end_row();

                ui.label("Organization");
                ui.text_edit_singleline(&mut award.organization);
                ui.end_row();

                ui.label("Date");
                ui.text_edit_singleline(&mut award.date);
                ui.end_row();
            });

        ui.add_space(10.0);
        editable_list(ui, "Points", TextEditMode::MultiLine, &mut award.points);
    });
}

pub fn project_ui(ui: &mut egui::Ui, project: &mut Project, i: usize) {
    egui::CollapsingHeader::new(if project.name.is_empty() {
        "New Project"
    } else {
        &project.name
    })
    .default_open(true)
    .id_salt(format!("project_{}", i))
    .show(ui, |ui| {
        egui::Grid::new(format!("project_grid_{}", i))
            .num_columns(2)
            .spacing([40.0, 4.0])
            .striped(true)
            .show(ui, |ui| {
                ui.label("Name");
                ui.text_edit_singleline(&mut project.name);
                ui.end_row();

                ui.label("GitHub URL");
                ui.text_edit_singleline(&mut project.github);
                ui.end_row();

                ui.label("GitHub Handle");
                ui.text_edit_singleline(&mut project.github_handle);
                ui.end_row();
            });

        ui.add_space(10.0);
        editable_list(
            ui,
            "CV Types",
            TextEditMode::SingleLine,
            &mut project.cv_type,
        );
        editable_list(ui, "Points", TextEditMode::MultiLine, &mut project.points);
    });
}

pub fn experience_ui(ui: &mut egui::Ui, experience: &mut Experience, i: usize) {
    egui::CollapsingHeader::new(format!(
        "{} ({})",
        if experience.role.is_empty() {
            "New Role"
        } else {
            &experience.role
        },
        if experience.company.is_empty() {
            "New Company"
        } else {
            &experience.company
        }
    ))
    .id_salt(format!("experience_{}", i))
    .show(ui, |ui| {
        egui::Grid::new(format!("experience_{}", i))
            .num_columns(2)
            .spacing([40.0, 4.0])
            .striped(true)
            .show(ui, |ui| {
                ui.label("Company");
                ui.text_edit_singleline(&mut experience.company);
                ui.end_row();

                ui.label("Location");
                ui.text_edit_singleline(&mut experience.location);
                ui.end_row();

                ui.label("Role");
                ui.text_edit_singleline(&mut experience.role);
                ui.end_row();

                ui.label("Dates");
                ui.text_edit_singleline(&mut experience.dates);
                ui.end_row();

                ui.label("Job Type");
                ui.text_edit_singleline(&mut experience.job_type);
                ui.end_row();
            });

        ui.add_space(10.0);
        editable_list(
            ui,
            "CV Types",
            TextEditMode::SingleLine,
            &mut experience.cv_type,
        );
        editable_list(
            ui,
            "Points",
            TextEditMode::MultiLine,
            &mut experience.points,
        );
    });
}

pub fn top_panel_ui(app: &App, ctx: &egui::Context, frame: egui::Frame) {
    egui::TopBottomPanel::top("top_panel")
        .frame(frame)
        .show(ctx, |ui| {
            ui.horizontal(|ui| {
                // --- Left side of the panel ---
                ui.with_layout(egui::Layout::left_to_right(egui::Align::Center), |ui| {
                    ui.label(egui::RichText::new("AutoCV Config").strong().size(24.0));
                });

                // --- Right side of the panel ---
                ui.with_layout(egui::Layout::right_to_left(egui::Align::Center), |ui| {
                    let compile_button = egui::Button::new(
                        egui::RichText::new("Compile PDF")
                            .color(crate::app::style::ayu_dark::BG_DARK)
                            .strong(),
                    )
                    .fill(crate::app::style::ayu_dark::ACCENT_YELLOW);

                    if ui
                        .add(compile_button)
                        .on_hover_text("Compile the CV configuration to a PDF file")
                        .clicked()
                    {
                        // TODO: Send compile message to actor
                        app.actor_sender
                            .try_send(ActorMessage::CompileCV(Box::new(app.local_state.clone())))
                            .unwrap_or_else(|e| {
                                log::error!("Failed to send CompileCV message: {}", e);
                            });
                        info!("CompileCV message sent to actor");
                    }
                });
            });
        });
}

pub fn side_panel_ui(ctx: &egui::Context, app: &mut App) {
    egui::SidePanel::left("side_panel")
        .min_width(150.0)
        .resizable(false)
        .show(ctx, |ui| {
            ui.add_space(10.0);

            // Center the buttons vertically and horizontally
            ui.with_layout(
                egui::Layout::top_down_justified(egui::Align::Center),
                |ui| {
                    ui.spacing_mut().item_spacing.y = 10.0;

                    // Use selectable_value for cleaner logic and consistent styling
                    let button_size = egui::vec2(140.0, 30.0);
                    if ui
                        .add_sized(
                            button_size,
                            egui::Button::selectable(
                                app.selected_tab == AppTab::General,
                                "General",
                            ),
                        )
                        .clicked()
                    {
                        app.selected_tab = AppTab::General;
                    }
                    if ui
                        .add_sized(
                            button_size,
                            egui::Button::selectable(
                                app.selected_tab == AppTab::Experiences,
                                "Experiences",
                            ),
                        )
                        .clicked()
                    {
                        app.selected_tab = AppTab::Experiences;
                    }
                    if ui
                        .add_sized(
                            button_size,
                            egui::Button::selectable(
                                app.selected_tab == AppTab::Projects,
                                "Projects",
                            ),
                        )
                        .clicked()
                    {
                        app.selected_tab = AppTab::Projects;
                    }
                },
            );
        });
}
