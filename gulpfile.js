const gulp = require('gulp');
const { log } = require('gulp-util');
const { exec } = require('child_process');
const path = require('path');
const fs = require('fs');
const { Transform } = require('stream');
require('colors');

let colorStreams = [];

class ColoredTransform extends Transform {
  constructor(procOut, color, { name, ...opts }) {
    super(opts);
    procOut.pipe(this);
    this.pipe(process.stdout);
    this.color = color || 'white';
    colorStreams.push(this);
  }

  _transform(chunk, encoding, callback) {
    this.push(chunk.toString()[this.color])
  }
}

const destroyColored = (name) => {
  colorStreams.forEach(stream => {
    if (stream.name === name) stream.destroy();
  });
  colorStreams = [];
};

const src = path.join(__dirname, 'src');
const glob = path.join(src, '*.go');
const out = path.resolve(src, 'a.out');

let goProc = null;
let compileSuccess = false;

gulp.task('compile', () => new Promise(resolve => {
  destroyColored('app');
  if (goProc) goProc.kill();
  compileSuccess = false;

  try {
    if (fs.statSync(out).isFile()) fs.unlinkSync(out);
  } catch (err) {
  }

  const proc = exec(`go build -o ${out} ${glob}`);

  new ColoredTransform(proc.stdout, 'blue', { name: 'compile' });
  new ColoredTransform(proc.stderr, 'red', { name: 'compile' });

  proc.on('close', (code) => {
    if (code > 0) log(`-- Compiler errored out with code ${code} --`.red);
    else compileSuccess = true;
    resolve()
  });
}));

gulp.task('run', () => new Promise(resolve => {
  destroyColored('compile');
  if (!compileSuccess) return resolve();

  goProc = exec(out);

  new ColoredTransform(goProc.stdout, 'red', { name: 'app' });
  new ColoredTransform(goProc.stderr, 'blue', { name: 'app' });

  goProc.on('close', (code) => {
    if (code > 0) log(`-- Go errored out with code ${code} --`.red);
    else log('-- Go finished smoothly --'.green);
    goProc = null;
  });

  resolve();
}));

gulp.task(
  'dev',
  () => gulp.watch(
    'src/**/*.go',
    { queue: true, ignoreInitial: false },
    gulp.series(['compile', 'run']),
  )
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
      p.on('close', () => resolve())
    }))
  );

  return new Promise(resolve => {
    const p = exec(`minica --domains localhost`);
    new ColoredTransform(p.stdout, 'yellow', { name: 'm' });
    new ColoredTransform(p.stderr, 'red', { name: 'm' });
    p.on('close', () => resolve())
  })
});