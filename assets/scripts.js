import htmx from 'htmx.org'

import Alpine from 'alpinejs'

// Add Alpine instance to window object.
window.Alpine = Alpine

// Start Alpine.
Alpine.start()

import {
  createIcons,
  replaceElement,
  Menu,
  ArrowRight,
  Globe,
  Moon,
  Bell,
  Circle,
  CircleCheck as CircleChecked,
  CalendarPlus as Schedule,
  RefreshCw as Refresh,
  Search,
  Sun,
  Trash2 as Trash,
  Edit,
  Plus,
  Minus,
  X,
  Cat,
  Rabbit,
  Snail,
  Squirrel,
  Turtle,
  Bird,
  Bug,
  Fish,
  Rat,
  Worm,
  CircleHelp as Unknown,
} from 'lucide';

import { 
  bearFace as Bear, 
  elephant as Elephant, 
  bullHead as Cow, 
  butterfly as Butterfly, 
  catBig as Lion, 
  crab as Crab,
  frogFace as Frog, 
  hedgehog as Hedgehog, 
  horseHead as Horse, 
  owl as Owl, 
  pig as Pig, 
  shark as Shark, 
  whale as Whale, 
  whaleNarwhal as Narwhal
 } from '@lucide/lab';

 const lucideConfig = {
  attrs: {
    class: ['icon'],
    'stroke-width': 1,
  },
  icons: {
    Unknown,
    Menu,
    ArrowRight,
    Globe,
    Moon,
    Bell,
    Circle,
    CircleChecked,
    Refresh,
    Schedule,
    Search,
    Sun,
    Trash,
    Edit,
    Plus,
    Minus,
    X,
    Cat,
    Rabbit,
    Snail,
    Squirrel,
    Turtle,
    Bird,
    Bug,
    Fish,
    Rat,
    Worm,
    Bear,
    Elephant,
    Cow,
    Butterfly,
    Lion,
    Crab,
    Frog,
    Hedgehog,
    Horse,
    Owl,
    Pig,
    Shark,
    Whale,
    Narwhal,
  }  
};

createIcons(lucideConfig);

function withMockDocument(mock, callback) {
  const originalDocument = document;

  // Redefine document for this scope
  global.document = mock;

  try {
    callback();
  } finally {
    global.document = originalDocument; // Restore the original document
  }
}

htmx.onLoad((content) => {
  console.log("htmx:onLoad", content);

  // Sub loaded content for document to confine ludide.createIcons()
  //  search for elements to replace to the new content rather than the whole document
  withMockDocument(content, function() {
    createIcons(lucideConfig);
  });
});