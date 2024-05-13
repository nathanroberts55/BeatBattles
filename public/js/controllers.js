// import "https://unpkg.com/@hotwired/turbo@8.0.4/dist/turbo.es2017-esm.js"
import { Application } from "https://unpkg.com/@hotwired/stimulus@3.2.2/dist/stimulus.js"
import SongController from "./controllers/song.js"

window.Stimulus = Application.start()
window.Stimulus.debug = window.location.protocol === "http:"
window.Stimulus.register("songs", SongController)


