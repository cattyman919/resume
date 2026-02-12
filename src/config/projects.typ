// Maybe create default points
// The points item can be override per each cv type instead

#let make-project(
  title: str,
  description: str,
  url:str,
  url_handle:str,
  points: ()
) = {
  return (
    title: title,
    description: description,
    url: url,
    url_handle: url_handle,
    points: points
  )
}


#let vxlang = make-project(
  title: "VxLang Security Analysis (Thesis)",
  description: "A comprehensive undergraduate thesis project analyzing the effectiveness of VxLang Code Virtualization.",
  url: "https://github.com/cattyman919/skripsi-project",
  url_handle:"cattyman919/skripsi-project",
  points: (
    [Investigated *data privacy protection methods* through *code virtualization* and obfuscation techniques, demonstrating a *>90% reduction* in symbol visibility against reverse engineering attacks.],
    [Conducted in-depth *static and dynamic analysis* using *Ghidra* and *x64dbg* to validate the integrity of protected binary data.],
    [Developed multiple C++ test applications (Console, *Qt*, *Dear ImGui*) and cryptographic benchmarks using *CMake*, *Ninja*, and *Clang*, integrating the *VxLang SDK* to implement granular virtualization on critical execution paths.]
  )
)

#let yazi = make-project(
  title: "Yazi (Open Source)",
  description: [Contributed to *Yazi*, a lightning-fast terminal file manager written in *Rust*, by implementing real-time progress visualization for I/O operations.],
  url: "https://github.com/sxyazi/yazi/pull/3121",
  url_handle: "sxyazi/yazi",
  points: (
    [Developed a *real-time task progression system* for file copying operations, replacing static status messages with a visual progress bar, byte transfer metrics, and file counts.],
    [Engineered the solution by modifying the core *Rust* task logic (`yazi-core`) to track I/O states and updating the *Lua*-based UI components (`yazi-plugin`) for rendering.],
    [Navigated a large, complex open-source codebase to implement the feature, successfully merging the contribution to enhance user feedback for long-running operations.]
  )
)

#let autocv = make-project(
  title: "AutoCV",
  description: "A dynamic CV generator that builds multiple tailored PDF resumes from a modular YAML data source using *Go* and *LaTeX*.",
  url: "https://github.com/cattyman919/resume",
  url_handle: "cattyman919/resume",
  points: (
        [Streamlined CV creation by developing a tool that generates *multiple customized PDF versions* from a modular data source split across three YAML files (`general`, `experiences`, `projects`), eliminating repetitive manual editing.],
    [Engineered a *concurrent build process* using *Go's goroutines and waitgroups*, significantly cutting down generation time for multiple CVs.],
    [Containerized the entire build environment with *Docker*, ensuring consistent and reproducible builds across different machines.]
  )
)

#let moviedb = make-project(
  title: "MovieDB Showcase",
  description: [A full-stack movie browsing application featuring a *Go* backend API that serves data from the TMDB API, and a modern *React* frontend.],
  url: "https://github.com/cattyman919/movies",
  url_handle: "cattyman919/movies",
  points: (
    [Architected a backend API using *Go* and the *Gin framework* to serve as a proxy for the TMDB API, handling requests for popular, top-rated, and upcoming movies.],
    [Developed a responsive frontend with *React*, *TypeScript*, and *Tailwind CSS* to display movie data fetched from the backend, creating an interactive user experience.],
    [Containerized the entire application (*Go backend, React frontend, MongoDB*) using *Docker* and *Docker Compose*, and wrote integration tests with *Testcontainers* for the database layer.],
  )
)

#let slash = make-project(
  title: "Slash",
  description:  "A 2D platformer game built from the ground up in *Go* using the *Ebitengine* game engine, featuring custom tilemap rendering, player/enemy mechanics, and a dynamic camera system.",
  url: "https://github.com/cattyman919/slash",
  url_handle: "cattyman919/slash",
  points: (
    [Developed a complete 2D game engine loop using *Go* and *Ebitengine*, managing game state, entity updates, and rendering.],
    [Integrated *Tiled* for level design by creating a custom parser to load and render tilemaps from JSON data, including collision detection logic.],
    [Implemented core gameplay mechanics including player movement (WASD), enemy AI that follows the player, and a dynamic camera that smoothly tracks the player's position.]
  )
)

#let restomatic = make-project(
  title: "RestoMatic",
  description: [A web platform for easy food and drink ordering. Users can browse restaurant menus, add funds, and leave ratings.],
  url:"https://github.com/SistemBasisData2023/RestoMatic",
  url_handle:"SistemBasisData2023/RestoMatic",
  points: (
    [Engineered a responsive and intuitive user interface with *React* and *Tailwind CSS*, leading the frontend development to enhance user engagement.],
    [Implemented a secure payment system with *frontend validation*, ensuring sufficient balance before processing transactions.],
  )
)

#let dancertos1 = make-project(
  title: "DanceRTOS",
  description:  [An attendance system using FreeRTOS ESP32 and RFID, with a web server and Blynk integration for class and schedule management.],
  url: "https://github.com/cattyman919/AbsenceSystem",
  url_handle: "cattyman919/AbsenceSystem",
  points: (
    [Led the end-to-end development of the DanceRTOS attendance system, building both the *Flutter frontend* and *NestJS backend*.],
    [Engineered a real-time student login system using *MQTT* and *RFID* on an *ESP32*, enabling secure and instantaneous attendance tracking.],
    [Designed and implemented a dynamic class schedule and attendance table, providing lecturers with an organized and *real-time view of student presence*.]
  )
)

#let dancertos2 = make-project(
  title: "DanceRTOS",
  description:  [An attendance system using FreeRTOS ESP32 and RFID, with a web server and Blynk integration for class and schedule management.],
  url: "https://github.com/cattyman919/AbsenceSystem",
  url_handle: "cattyman919/AbsenceSystem",
  points: (
    [Designed and implemented a dynamic class schedule and attendance table, providing lecturers with an organized and *real-time view of student presence*.],
    [Led the end-to-end development of the DanceRTOS attendance system, building both the *Flutter frontend* and *NestJS backend*.],
    [Engineered a real-time student login system using *MQTT* and *RFID* on an *ESP32*, enabling secure and instantaneous attendance tracking.]
  )
)

#let jaga = make-project(
  title: "Jaga",
  description: [A vehicle maintenance app with GPS tracking to monitor kilometers and notify users of service needs based on distance and time.],
  url:"https://github.com/cattyman919/Jaga",
  url_handle: "cattyman919/Jaga",
  points: (
    [Led the full-stack development of the Jaga vehicle maintenance app, creating a seamless user experience with *Flutter* and *NestJS*.],
[Developed a *GPS-based tracking system* to monitor vehicle mileage and trigger timely service reminders.],
    [Implemented a *multi-vehicle management system*, allowing users to track and manage maintenance for their entire fleet.]
  )
)

#let jsleep = make-project(
  title: "JSleep",
  description: [A full-stack room reservation backend built with Java and Spring Boot, featuring a custom JSON-based database engine for managing users, rooms, and bookings.],
  url: "https://github.com/cattyman919/JSleep",
  url_handle: "cattyman919/JSleep",
  points:(
    [Developed a complete RESTful API using *Java* and *Spring Boot* to manage the entire hotel booking lifecycle, from user authentication to room reservations.],
    [Engineered a custom, lightweight database solution using a *JSON-based file system (`JsonDBEngine`)*, enabling persistent data storage and retrieval without a traditional database server.],
    [Implemented core business logic for dynamic room filtering, availability checking based on booking dates, and a complete payment processing system including balance management and voucher application.]
  )
)

#let electronic-vault-lock = make-project(
  title: "Electronic Vault Lock",
  description: [A secure 4-digit combination lock for protecting items, offering a simple and reliable locking mechanism.],
  url: "https://github.com/rroiii/Electronic-Vault-Lock",
  url_handle: "rroiii/Electronic-Vault-Lock",
  points:(
    [Authored and optimized the *VHDL code* for the electronic vault lock, ensuring robust and reliable security.],
    [Conducted extensive simulations in *ModelSim* to validate the system's performance and identify potential vulnerabilities.],
    [Led the synthesis of the design in *Quartus*, optimizing for area and power to create an efficient and compact solution.]
  )
)

#let vaio = make-project(
  title: "VAIO (Vacuum All in One)",
  description: [VAIO is an ESP32-based robot featuring autonomous navigation, a vacuum system, and multiple intuitive control methods including a sensor-equipped glove, voice commands via speech recognition, and a standard PS4 controller.],
  url: "https://github.com/VAIO-CE/VAIO-Code",
  url_handle: "VAIO-CE/VAIO-Code",
  points: (
    [Architected a modular and scalable codebase in *C++* on the *PlatformIO IDE*, enabling efficient team collaboration and future development.],
    [Engineered a real-time gyroscope-based control system with *ESP-NOW*, allowing for intuitive and responsive robot control via a wearable glove.],
    [Implemented a multi-threaded *FreeRTOS* environment to manage concurrent tasks, ensuring seamless switching between autonomous, gesture, and PS4 control modes.]
  )
)

#let http-server = make-project(
  title:"HTTP Server",
  description: [A custom HTTP server written in C that leverages low-level socket programming and polling to handle multiple client connections. It features dynamic routing via a binary search tree and simple HTTP parsing for flexible response handling.],
  url: "https://github.com/cattyman919/http",
  url_handle: "cattyman919/http",
  points: (
    [Engineered a *non-blocking, polling-based socket server* in *C* to handle multiple concurrent client connections efficiently.],
    [Implemented a *custom HTTP parser and response generator*, enabling the server to handle a variety of requests and serve dynamic content.],
    [Designed and implemented a *binary search tree for dynamic routing*, allowing for efficient and scalable URL-to-file mapping.]
  )
)


#let home-server1 = make-project(
  title: "Home Server",
  description: [Dell OptiPlex 7050 running Debian as the Operating System used for my personal custom server. It currently provides services such as Portainer, Jellyfin, CasaOS, Samba Server, SIP Server, and also used as a backup for my files. Connected remotely using Tailscale.],
  url: none,
  url_handle : none,
  points: (
    [*Linux Server (Debian)* that deploys and manage a suite of services including *Jellyfin for media streaming*, *Portainer for Docker container management*, and an *Asterisk SIP server for home VoIP*.],
    [Configured core network services, including *Samba for local file sharing*, *nginx as a reverse proxy* for simplified service access, and implemented *fail2ban and firewalld* for enhanced security.]
  )
)

#let home-server2 = make-project(
  title: "Home Server",
  description: [Dell OptiPlex 7050 running Debian as the Operating System used for my personal custom server. It currently provides services such as Portainer, Jellyfin, CasaOS, Samba Server, SIP Server, and also used as a backup for my files. Connected remotely using Tailscale.],
  url: none,
  url_handle : none,
  points: (
    [Deployed and manage a suite of services including *Jellyfin for media streaming*, *Portainer for Docker container management*, and an *Asterisk SIP server for home VoIP*.],
    [Configured core network services, including *Samba for local file sharing*, *nginx as a reverse proxy* for simplified service access, and implemented *fail2ban and firewalld* for enhanced security.],
    [Established a robust data management and remote access strategy using *rsync for automated backups* and *Tailscale* for secure, seamless access to the entire home network from any location.]
  )
)
