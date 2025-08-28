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

// Global function to access Alpine.js $data from any element
window.alpine = {
    data: (field) => {
        const element = document.querySelector('main[x-data]');
        if (element && element._x_dataStack && element._x_dataStack.length > 0) {
            if (field) {
                return element._x_dataStack[0][field];

            } else {
                return element._x_dataStack[0];
            }
        }

        return null;
    }
};
