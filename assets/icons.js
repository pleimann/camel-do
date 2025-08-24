import htmx from 'htmx.org';

import {
  createIcons,
  Menu,
  ArrowRight,
  Globe,
  Moon,
  Bell,
  Circle,
  CircleCheck as CircleChecked,
  CalendarPlus as ScheduleNow,
  CalendarMinus as Backlog,
  CalendarClock as ScheduleTime,
  RefreshCw as Refresh,
  Search,
  Pencil,
  PencilLine,
  Sun,
  Trash2 as Trash,
  Edit,
  Plus,
  Minus,
  X,
  Package,
  Boxes as Packages,
  PackageOpen,
  PackagePlus,
  PackageMinus,
  NotepadText,
  Clock,
  CircleHelp as Unknown,
  ChevronDown,
  ChevronUp,
  ChevronLeft,
  ChevronRight,

  Bird,
  Bug,
  Cat,
  Dog,
  Fish,
  Panda,
  Rabbit,
  Rat,
  Snail,
  Squirrel,
  Turtle,
  Worm,
} from 'lucide';

import { 
  bearFace as Bear, 
  bee as Bee,
  cowHead as Cow, 
  butterfly as Butterfly, 
  catBig as Lion, 
  crab as Crab,
  elephant as Elephant, 
  frogFace as Frog, 
  hedgehog as Hedgehog, 
  horseHead as Horse, 
  owl as Owl, 
  pig as Pig, 
  shark as Shark, 
  spider as Spider,
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
    ChevronDown,
    ChevronUp,
    ChevronLeft,
    ChevronRight,
    Globe,
    Moon,
    Bell,
    Circle,
    CircleChecked,
    Refresh,
    ScheduleNow,
    ScheduleTime,
    Backlog,
    Pencil,
    PencilLine,
    Search,
    Sun,
    Trash,
    Edit,
    Plus,
    Minus,
    X,
    Package,
    Packages,
    PackageOpen,
    PackagePlus,
    PackageMinus,
    NotepadText,
    Clock,
    
    Bear,
    Bee,
    Bird,
    Bug,
    Butterfly,
    Cat,
    Crab,
    Cow,
    Dog,
    Elephant,
    Fish,
    Frog,
    Hedgehog,
    Horse,
    Lion,
    Narwhal,
    Owl,
    Panda,
    Pig,
    Rabbit,
    Rat,
    Snail,
    Squirrel,
    Turtle,
    Worm,
    Shark,
    Spider,
    Whale,
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
  // Sub loaded content for document to confine ludide.createIcons()
  //  search for elements to replace to the new content rather than the whole document
  withMockDocument(content, function() {
    createIcons(lucideConfig);
  });
});
