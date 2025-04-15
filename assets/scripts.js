import htmx from 'htmx.org';

window.htmx = htmx;

import 'htmx-ext-response-targets'; 
import 'htmx-ext-remove-me'; 
import './hx-drag';

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

Alpine.directive('log', (el, { expression }, { evaluateLater, effect }) => {
    let getThingToLog = evaluateLater(expression)
 
    effect(() => {
        getThingToLog(thingToLog => {
            console.log(thingToLog)
        })
    })
})
