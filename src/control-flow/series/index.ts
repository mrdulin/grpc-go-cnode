import async from 'async';
import request from 'request-promise';

import { config } from '../../config';

function getBookByName(name: string) {
  const url = `${config.bookstore.apiUrl}/search/${name}`;
  const options = {
    uri: url,
    method: 'GET',
    json: true
  };
  return request(options);
}

function getBooksSeries(names: string[]) {
  const start = Date.now();
  return new Promise((resolve, reject) => {
    const tasks = names.map((name: string) => {
      return function task(callback) {
        getBookByName(name)
          .then(response => {
            if (response.error !== '0') {
              callback(new Error(`get books of ${name} failed.`));
            }
            callback(null, response);
          })
          .catch(callback);
      };
    });
    async.series(tasks, function(error, responses) {
      console.log(`completed in ${Date.now() - start} ms`);
      if (error) {
        reject(error);
      }
      resolve(responses);
    });
  });
}

function getStarted() {
  const start = Date.now();
  async.series(
    {
      one(callback) {
        setTimeout(function() {
          callback(null, 1);
        }, 1000);
      },
      two(callback) {
        setTimeout(function() {
          callback(null, 2);
        }, 1000);
      }
    },
    function(err, results) {
      const end = Date.now();
      console.log(`completed in ${end - start}ms`);
      // results is now equal to: {one: 1, two: 2}
    }
  );
}

export { getBooksSeries, getStarted };
