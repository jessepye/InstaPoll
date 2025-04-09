const fs = require('fs');
const path = require('path');

function test() {
  const requiredFiles = [
    'index.html',
    'error.html'
  ];

  const missingFiles = requiredFiles.filter(file => {
    return !fs.existsSync(path.join(__dirname, file));
  });

  if (missingFiles.length > 0) {
    console.error('Missing required files:', missingFiles);
    process.exit(1);
  }

  console.log('All required files exist!');
}

test(); 