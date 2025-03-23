import htmx from 'htmx.org';

window.htmx = htmx;

import 'htmx-ext-response-targets'; 
import 'htmx-ext-remove-me'; 
import './hx-drag';

import './icons'

import Alpine from 'alpinejs'

// Add Alpine instance to window object.
window.Alpine = Alpine

// Start Alpine.
Alpine.start()
