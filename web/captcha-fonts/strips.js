const { execSync } = require('child_process');
const fs = require('fs');
const { resolve, basename } = require('path');

const projectDir = resolve(__dirname, '..');
const rootDir = resolve(projectDir, '..');

(async () => {
  execSync(`rm -rf ${__dirname}/striped/*.ttf`);

  const fonts = execSync(`find ${__dirname}/source -type f -name "*.ttf"`, {
    encoding: 'utf8'
  })
    .trim()
    .split('\n');

  fonts.forEach((file) => {
    const f = basename(file);
    execSync(
      `pyftsubset ${file} --name-IDs="" --unicodes=U+0030-0039 --output-file=${__dirname}/striped/${f}`
    );
  });

  const striped = execSync(`find ${__dirname}/striped -type f -name "*.ttf"`, {
    encoding: 'utf8'
  })
    .trim()
    .split('\n');

  const GoModuleFonts = [];
  GoModuleFonts.push(`package main

import (
  "encoding/base64"
)

// CaptchaFonts list of striped fonts for captchaUsage
var CaptchaFonts map[int][]byte
// FarsiCaptchaFonts list of striped fonts for captchaUsage
var FarsiCaptchaFonts map[int][]byte

func init() {
  CaptchaFonts = make(map[int][]byte)
  FarsiCaptchaFonts = make(map[int][]byte)`);

  let farsiCount = 0;
  let otherCount = 0;
  striped.forEach((file) => {
    const b64 = fs.readFileSync(file, { encoding: 'base64' });
    if (file.match(/fa_/)) {
      GoModuleFonts.push(
        `  FarsiCaptchaFonts[${farsiCount}], _ = base64.StdEncoding.DecodeString("${b64}")`
      );
      farsiCount += 1;
    } else {
      GoModuleFonts.push(
        `  CaptchaFonts[${otherCount}], _ = base64.StdEncoding.DecodeString("${b64}")`
      );
      otherCount += 1;
    }
  });
  GoModuleFonts.push('}');

  fs.writeFileSync(`${rootDir}/captcha_font.go`, GoModuleFonts.join('\n'), {
    encoding: 'utf8'
  });
})();
