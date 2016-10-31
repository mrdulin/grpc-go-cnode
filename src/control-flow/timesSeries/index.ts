import timesSeries from 'async/timesSeries';

const interval = 5 * 1000;
const scheduler = 10 * 1000;
const executionTimes = scheduler / interval;

function request(operation: string) {
  return new Promise(resolve => {
    setTimeout(() => {
      console.log(operation);
      resolve();
    }, 1000);
  });
}

function sleep(ms: number) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

async function createUser(id, callback) {
  await request(`create user with id:${id} - timestamp: ${Date.now()}`);
  await sleep(interval);
  callback();
}

function main() {
  console.log('main');
  timesSeries(executionTimes, createUser, (err, users) => {
    if (err) {
      console.error(err);
      return;
    }
    console.log('create users done.');
  });
}

main();
