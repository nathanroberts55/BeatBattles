import { Application } from "https://unpkg.com/@hotwired/stimulus@3.2.2/dist/stimulus.js"
import SongController from "./controllers/song.js"
import FormController from "./controllers/form.js"

window.Stimulus = Application.start()
window.Stimulus.debug = window.location.protocol === "http:"
window.Stimulus.register("song", SongController)
window.Stimulus.register("form", FormController)

