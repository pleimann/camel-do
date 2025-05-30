import "./htmx"

import 'htmx-ext-response-targets'; 
import 'htmx-ext-remove-me'; 
import 'htmx-ext-debug';
import 'hx-drag';

import './icons'

import Alpine from 'alpinejs'
import focus from '@alpinejs/focus'
import anchor from '@alpinejs/anchor'

// Add Alpine instance to window object.
window.Alpine = Alpine

// Start Alpine.
Alpine.plugin(focus)
Alpine.plugin(anchor)
Alpine.start()
