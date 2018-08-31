const detect = require('async/detect');
const fs = require('fs');
const path = require('path');

const files = [];
for (let i = 0; i < 3; i++) {
  files.push(path.resolve(__dirname, `file${i}.txt`));
}

const fsAccessCheck = (filePath, callback) => {
  fs.access(filePath, fs.constants.F_OK, err => {
    callback(null, !err);
  });
};

/**
 * 返回files中第一个\通过异步函数fsAccessCheck测试的值
 */
detect(files, fsAccessCheck, (err, result) => {
  if (err) return console.error(err.stack);
  console.log(result); //结果：  /Users/dulin/workspace/Training.nodejs/async/collections/detect/file2.txt
});
