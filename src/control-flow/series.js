const async = require('async');

const start = new Date().getTime();

async.series(
  [
    function(callback) {
      setTimeout(() => {
        callback(null, 'angular');
      }, 100);
    },
    function(callback) {
      setTimeout(() => {
        callback(null, 'react');
      }, 200);
    },
    function(callback) {
      setTimeout(() => {
        callback(null, 'rxjs');
      }, 300);
    }
  ],
  function(err, results) {
    console.log('completed in ', new Date().getTime() - start + 'ms');
    console.log(results);
  }
);
