use autocv_core::cv_model::{
    Award, Certificate, Education, Experience, GeneralCVData, PersonalInfo, Project,
    SkillsAchievements,
};

use crate::app::{App, AppTab};

// A helper function to create a consistent section header
fn section_header(ui: &mut egui::Ui, title: &str) {
    ui.add_space(10.0);
    ui.heading(title);
    ui.add_space(4.0);
}

// A helper for adding/removing items from a Vec<String>
fn editable_list(ui: &mut egui::Ui, title: &str, items: &mut Vec<String>) {
    ui.label(title);
    let mut item_to_remove = None;
    for (i, item) in items.iter_mut().enumerate() {
        ui.horizontal(|ui| {
            ui.add(egui::TextEdit::multiline(item).desired_width(f32::INFINITY));
            // Use a styled button for better looks
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
            ui.add(egui::TextEdit::singleline(&mut cert.name).desired_width(f32::INFINITY));
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

pub fn general_ui(ui: &mut egui::Ui, general_cv: &mut GeneralCVData) {
    personal_info_ui(ui, &mut general_cv.personal_info);
    skills_achivements_ui(ui, &mut general_cv.skills_achievements);

    section_header(ui, "Education");
    // TODO: Add a button to add/remove education entries
    for education in general_cv.education.iter_mut() {
        education_ui(ui, education);
    }

    section_header(ui, "Awards");
    // TODO: Add a button to add/remove award entries
    for award in general_cv.awards.iter_mut() {
        awards_ui(ui, award);
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
        editable_list(ui, "Hard Skills", &mut skills_achievements.hard_skills);
        editable_list(ui, "Soft Skills", &mut skills_achievements.soft_skills);
        editable_list(
            ui,
            "Programming Languages",
            &mut skills_achievements.programming_languages,
        );
        editable_list(ui, "Databases", &mut skills_achievements.databases);
        editable_list(ui, "Misc Tools", &mut skills_achievements.misc);
        certificate_list(ui, "Certificates", &mut skills_achievements.certificates);
    });
}

pub fn education_ui(ui: &mut egui::Ui, education: &mut Education) {
    egui::CollapsingHeader::new(if education.institution.is_empty() {
        "New Entry"
    } else {
        &education.institution
    })
    .show(ui, |ui| {
        egui::Grid::new(format!("education_grid_{}", education.institution))
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
        editable_list(ui, "Details", &mut education.details);
    });
}

pub fn awards_ui(ui: &mut egui::Ui, award: &mut Award) {
    egui::CollapsingHeader::new(if award.title.is_empty() {
        "New Award".to_string()
    } else {
        award.title.clone()
    })
    .show(ui, |ui| {
        egui::Grid::new(format!("award_grid_{}", award.title))
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
        editable_list(ui, "Points", &mut award.points);
    });
}

pub fn project_ui(ui: &mut egui::Ui, project: &mut Project) {
    egui::CollapsingHeader::new(if project.name.is_empty() {
        "New Project"
    } else {
        &project.name
    })
    .show(ui, |ui| {
        egui::Grid::new(format!("project_grid_{}", project.name))
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
            });

        ui.add_space(10.0);
        editable_list(ui, "Points", &mut project.points);
        editable_list(ui, "CV Types", &mut project.cv_type);
    });
}

pub fn experience_ui(ui: &mut egui::Ui, experience: &mut Experience) {
    egui::CollapsingHeader::new(format!(
        "{} at {}",
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
    .show(ui, |ui| {
        egui::Grid::new(format!(
            "experience_grid_{}_{}",
            experience.company, experience.role
        ))
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
        editable_list(ui, "Points", &mut experience.points);
        editable_list(ui, "CV Types", &mut experience.cv_type);
    });
}

pub fn side_panel_ui(ctx: &egui::Context, app: &mut App) {
    egui::SidePanel::left("side_panel")
        .min_width(150.0)
        .resizable(false)
        .show(ctx, |ui| {
            ui.add_space(10.0);
            ui.heading("Config Tab");
            ui.add_space(5.0);
            ui.separator();
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

            // ui.spacing_mut().item_spacing.y = 10.0;
            // // Get the style once to reuse it
            // let theme = &ctx.style().visuals.clone();
            // let selected_fill = theme.selection.bg_fill;
            // let selected_stroke = theme.selection.stroke;
            //
            // // -- General Tab Button --
            // let is_general_selected = app.selected_tab == AppTab::General;
            // let general_button = egui::Button::new("General")
            //     .min_size(egui::vec2(120.0, 25.0)) // Give it a fixed size
            //     .fill(if is_general_selected {
            //         selected_fill
            //     } else {
            //         egui::Color32::TRANSPARENT
            //     })
            //     .stroke(if is_general_selected {
            //         selected_stroke
            //     } else {
            //         egui::Stroke::NONE
            //     });
            //
            // if ui.add(general_button).clicked() {
            //     app.selected_tab = AppTab::General;
            // }
            //
            // // -- Experiences Tab Button --
            // let is_experiences_selected = app.selected_tab == AppTab::Experiences;
            // let experiences_button = egui::Button::new("Experiences")
            //     .min_size(egui::vec2(120.0, 25.0))
            //     .fill(if is_experiences_selected {
            //         selected_fill
            //     } else {
            //         egui::Color32::TRANSPARENT
            //     })
            //     .stroke(if is_experiences_selected {
            //         selected_stroke
            //     } else {
            //         egui::Stroke::NONE
            //     });
            //
            // if ui.add(experiences_button).clicked() {
            //     app.selected_tab = AppTab::Experiences;
            // }
            //
            // // -- Projects Tab Button --
            // let is_projects_selected = app.selected_tab == AppTab::Projects;
            // let projects_button = egui::Button::new("Projects")
            //     .min_size(egui::vec2(120.0, 25.0))
            //     .fill(if is_projects_selected {
            //         selected_fill
            //     } else {
            //         egui::Color32::TRANSPARENT
            //     })
            //     .stroke(if is_projects_selected {
            //         selected_stroke
            //     } else {
            //         egui::Stroke::NONE
            //     });
            //
            // if ui.add(projects_button).clicked() {
            //     app.selected_tab = AppTab::Projects;
            // }
        });
}
