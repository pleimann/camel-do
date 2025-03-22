import htmx from 'htmx.org'

// htmx.logAll()

import Alpine from 'alpinejs'

// Add Alpine instance to window object.
window.Alpine = Alpine

// Start Alpine.
Alpine.start()

import {
  createIcons,
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

createIcons({
  attrs: {
    class: ['icon'],
    'stroke-width': 1,
  },
  nameAttr: 'data-lucide',
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
});
