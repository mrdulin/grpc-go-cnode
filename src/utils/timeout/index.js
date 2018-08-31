const timeout =  require('async/timeout');
const superagent = require('superagent');
const fs = require('fs');
const path = require('path');

//需求：异步操作的时间小于1秒时，正常执行，如果异步操作的时间大于1秒，抛出err对象，'ETIMEDOUT'

const searchBook = (query, cb) => {
	const url = `http://it-ebooks-api.info/v1/search/${query}`;

	superagent.get(url).then(res => {
		cb(null, {query, body: res.body});
	}).catch(err => {
		cb(err);
	});
}

const searchBookWrapped = timeout(searchBook, 800);
const query = 'angular';
const retryDelay = 1000, maxRetryCount = 5;
searchBookWrapped.retryCount = 1;

const searchHandle = (err, data) => {
	if(err) {
		if(err.code === 'ETIMEDOUT') {
			const {retryCount} = searchBookWrapped;
			if(retryCount < maxRetryCount) {
				console.log(`请求超时，正在重新尝试请求. ${retryCount}次`)
				setTimeout(() => {
					searchBookWrapped.retryCount++;
					searchBookWrapped.call(null, query, searchHandle);
					return;
				}, retryDelay);
			} else {
				return console.error(err.stack);
			}
		} else {
			return console.error(err.stack);
		}
	} else {
		const result = {data};
		fs.writeFileSync(path.resolve(__dirname, 'data.json'), JSON.stringify(result, null, 4));
		console.log('请求成功，数据写入完毕!');
	}
}

searchBookWrapped(query, searchHandle);

