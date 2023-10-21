import 'regenerator-runtime/runtime'

import { Load } from './general';

window.addEventListener('load', () => {
  Load()
});

if (window.history && window.history.pushState) {
  window.history.pushState = new Proxy(window.history.pushState, {
    apply: (target, thisArg, argArray) => {
      Load()
      return target.apply(thisArg, argArray);
    },
  });  
}
