import { library, config } from '@fortawesome/fontawesome-svg-core'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import {
	faSearch,
	faGamepad,
	faVideo,
	faChevronLeft,
	faChevronRight,
	faExclamationCircle,
	faSpinner,
	faPen,
	faTimes,
	faCheckCircle,
	faTrash
} from "@fortawesome/free-solid-svg-icons";

import { faYoutube } from "@fortawesome/free-brands-svg-icons";

const icons = [
	faSearch,
	faGamepad,
	faVideo,
	faChevronLeft,
	faChevronRight,
	faExclamationCircle,
	faSpinner,
	faPen,
	faYoutube,
	faTimes,
	faCheckCircle,
	faTrash,
];

// This is important, we are going to let Nuxt worry about the CSS
config.autoAddCss = false;

// You can add your icons directly in this plugin. See other examples for how you
// can add other styles or just individual icons.
library.add(...icons);

export default defineNuxtPlugin((nuxtApp) => {
	nuxtApp.vueApp.component('FontAwesomeIcon', FontAwesomeIcon)
})
