const gulp = require('gulp');
const { log } = require('gulp-util');
const { exec, execFile } = require('child_process');
const path = require('path');
const fs = require('fs');
const { Transform } = require('stream');
require('colors');

const src = path.join(__dirname, 'src');
const glob = path.join(src, '*.go');
const out = path.resolve(src, 'a.out');

let colorStreams = [];
let redis = null;
let goProc = null;
let compileSuccess = false;
let retry = 0;

class ColoredTransform extends Transform {
  constructor(procOut, color, { name, ...opts }) {
    super(opts);
    procOut.pipe(this);
    this.pipe(process.stdout);
    this.color = color || 'white';
    this.procOut = procOut;
    colorStreams.push(this);
  }

  _transform(chunk, encoding, callback) {
    this.push(chunk.toString()[this.color])
  }

  stop() {
    this.emit('end')
    this.unpipe(process.stdout)
    procOut.unpipe(this)
    this.destroy()
  }
}

const destroyColored = (name) => {
  colorStreams.forEach((stream, i) => {
    if (stream.name === name) {
      stream.stop();
    }
    colorStreams.splice(i, 1)
  });
};

const murderRedis = () => new Promise(resolve => {
  const serialKiller = exec('docker container rm -f redis-authorizer');
  new ColoredTransform(redis.stdout, 'yellow', { name: 'redis' })
  new ColoredTransform(redis.stderr, 'red', { name: 'redis' })
  serialKiller.on('close', (code) => {
    if (code > 0) {
      util.log('FAILED TO KILL REDIS CONTAINER'.red);
      process.exit(code);
    }
    resolve();
  })
});

const startRedis = async () => {
  console.log('--- Starting Redis Container ---\n'.green)
  const redis = exec('docker run --name redis-authorizer --rm -p 6379:6379 redis');

  new ColoredTransform(redis.stdout, 'yellow', { name: 'redis' })
  new ColoredTransform(redis.stderr, 'red', { name: 'redis' })

  return redis.on('close', (code => {
    if (code > 0) {
      if (code === 125 && retry < 10) {
        retry++;
        murderRedis()
          .then(() => startRedis);
        return;
      }
      return util.log(`--- Redis exited with code ${code} :'( ---\n\n`.red);
      process.exit(code);
    }
    return util.log('--- Redis exited cleanly! ---'.green);
  }))
};


gulp.task('compile', () => new Promise(resolve => {
  destroyColored('app');
  if (goProc) goProc.kill();
  compileSuccess = false;

  try {
    if (fs.statSync(out).isFile()) fs.unlinkSync(out);
  } catch (err) {
  }

  const p = exec(`go build -o ${out} ${glob}`);

  new ColoredTransform(p.stdout, 'blue', { name: 'compile' });
  new ColoredTransform(p.stderr, 'red', { name: 'compile' });

  p.on('close', (code) => {
    if (code > 0) log(`-- Compiler errored out with code ${code} --`.red);
    else compileSuccess = true;
    resolve()
  });
}));

gulp.task('run', () => new Promise(resolve => {
  destroyColored('compile');
  if (!compileSuccess) return resolve();

  goProc = execFile(out,);

  // TODO: figure out why color stream drops logs when running go programs
  goProc.stdout.on("data",  chunk => process.stdout.write(chunk.toString().blue));
  goProc.stdout.on("error",  chunk => process.stdout.write(chunk.toString().red));
  goProc.stderr.on("data",  chunk => process.stdout.write(chunk.toString().red));
  goProc.stderr.on("error", chunk => process.stdout.write(chunk.toString().red));

  goProc.on('close', (code) => {
    if (code > 0) log(`-- Go errored out with code ${code} --`.red);
    else log('-- Go finished smoothly --'.green);
    goProc = null;
  });

  resolve();
}));

gulp.task(
  'dev',
  async () => {
    await startRedis()
    return gulp.watch(
      'src/**/*.go',
      { queue: true, ignoreInitial: false },
      gulp.series(['compile', 'run']),
    )
  }
);

gulp.task("install", async () => {
  let deps = fs
    .readFileSync(path.join(__dirname, 'deps.txt'))
    .toString('UTF8')
    .split('\n')
    .filter(item => item.length > 0);

  log(
    '\n-- Installing --\n'.yellow,
    deps.map(dep => dep.yellow).join('\n'),
    '\n'
  );

  await Promise.all(
    deps.map((dep) => new Promise(resolve => {
      const p = exec(`go get -t ${dep}`);
      new ColoredTransform(p.stdout, 'yellow', { name: dep });
      new ColoredTransform(p.stderr, 'red', { name: dep });
      p.on('close', () => {
        destroyColored(dep);
        resolve()
      })
    }))
  );

  let images = ['redis']

  log(
    '\n-- Pulling --\n'.yellow,
    images.map(image => image.yellow).join('\n'),
    '\n'
  );

  await Promise.all(
    images.map((image) => new Promise(resolve => {
      const p = exec(`docker pull ${image}`);
      new ColoredTransform(p.stdout, 'yellow', { name: image });
      new ColoredTransform(p.stderr, 'red', { name: image });
      p.on('close', () => {
        destroyColored(image)
        resolve()
      })
    }))
  );

  log('\n-- Creating development ssl cert --\n'.yellow)

  let l = false;

  try {
    l = fs.statSync(path.join(__dirname, 'localhost')).isDirectory();
  } catch (err) {
  } finally {
    [
      ...fs.readdirSync(__dirname).map(it => path.join(__dirname, it)),
      ...(
        l ?
          fs.readdirSync(path.join(__dirname, 'localhost'))
            .map(it => path.join(__dirname, 'localhost', it))
          : []
      )
    ].forEach(f => {
      if (f.split('.').pop() === 'pem')
        fs.unlinkSync(f);
    });
  }


  return new Promise(resolve => {
    const p = exec(`minica --domains localhost`);
    new ColoredTransform(p.stdout, 'yellow', { name: 'm' });
    new ColoredTransform(p.stderr, 'red', { name: 'm' });
    p.on('close', () => resolve())
  })
});