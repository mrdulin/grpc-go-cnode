const concat = require('async/concat');
const superagent = require('superagent');
const fs = require('fs');
const path = require('path');
const stream = require('stream');

function isStream(obj) {
  return obj instanceof stream.Stream;
}

function isReadable(obj) {
  return isStream(obj) && typeof obj._read == 'function' && typeof obj._readableState == 'object';
}

function isWritable(obj) {
  return isStream(obj) && typeof obj._write == 'function' && typeof obj._writableState == 'object';
}

const rUrl = 'http://it-ebooks-api.info/v1/';
let pathNames = ['search/angular', 'search/react', 'search/node', 'search/jquery', 'search/backbone'];

const getQuery = url => {
  const urlPartial = url.split('/');
  const query = urlPartial[urlPartial.length - 1];
  return query;
};

let urlArr = pathNames.map(pathname => (item = rUrl + pathname));
let urlMap = pathNames.map(pathname => {
  const query = getQuery(pathname);
  return { [query]: rUrl + pathname };
});

const fetchData = (item, cb) => {
  console.log(item);
  let query = '',
    url = '';
  if (typeof item === 'object') {
    query = Object.keys(item)[0];
    url = item[query];
  } else {
    query = getQuery(item);
    url = item;
  }
  console.log(`${query}--数据请求中...`);

  //promise写法
  return superagent
    .get(url)
    .then(res => {
      // console.log('res readable: ', isReadable(res));
      // console.log('res writeable: ', isWritable(res));
      cb(null, { query, body: res.body });
    })
    .catch(err => {
      cb(err.response.error);
    });

  //callback写法
  // return superagent.get(url).end((err, res) => {
  // 	if(err) {
  // 		cb(err.response.error);
  // 	} else {
  // 		cb(null, {query, body: res.body});
  // 	}
  // });
};

/**
 * 集合是Array类型
 * 返回的结果顺序不固定，不一定和遍历的集合顺序一致
 * fetchData是并行执行的
 */
// concat(urlArr, fetchData, (err, results) => {
// 	if(err) return console.error(err.message);
// 	// console.log(results);
// 	const data = {data: results};
// 	console.log('开始写入数据...')
// 	fs.writeFileSync(path.resolve(__dirname, 'data.json'), JSON.stringify(data, null, 4), 'utf-8');
// 	console.log('数据写入完毕!')
// });

/**
 *	如果是对象类型，每个item的结构{ angular: 'http://it-ebooks-api.info/v1/search/angular' }
 */
concat(urlMap, fetchData, (err, results) => {
  if (err) return console.error(err.message);
  // console.log(results);
  const data = { data: results };
  // fs.writeFile('data.json', JSON.stringify(data, null, 4), err => {
  // 	if(err) throw err;
  // 	console.log('数据写入完毕!');
  // })

  const ws = fs.createWriteStream(path.resolve(__dirname, 'data.json'));
  ws.write(JSON.stringify(data, null, 4), 'utf8');
  ws.end();
  ws.on('finish', () => {
    console.log('数据写入完毕!');
  }).on('error', err => {
    console.error(err.stack);
  });
});
