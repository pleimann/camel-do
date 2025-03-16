import 'htmx.org'

import Alpine from 'alpinejs'
import { createIcons, icons } from 'lucide';

// Add Alpine instance to window object.
window.Alpine = Alpine

// Start Alpine.
Alpine.start()

// Recommended way, to include only the icons you need.
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
  Edit,
  Plus,
  Minus,
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
  icons: {
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
