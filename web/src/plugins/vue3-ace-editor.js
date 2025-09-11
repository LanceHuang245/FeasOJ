import { VAceEditor } from 'vue3-ace-editor';
import 'ace-builds/src-noconflict/mode-json';
import modeJava from "ace-builds/src-noconflict/mode-java?url";
import modeCpp from "ace-builds/src-noconflict/mode-c_cpp?url";
import modePy from "ace-builds/src-noconflict/mode-python?url";
import modeRust from "ace-builds/src-noconflict/mode-rust?url";
import modePHP from "ace-builds/src-noconflict/mode-php?url";
import modePascal from "ace-builds/src-noconflict/mode-pascal?url";
import { config } from "ace-builds";

// Theme
import 'ace-builds/src-noconflict/theme-chrome';
import "ace-builds/src-noconflict/theme-github";
import "ace-builds/src-noconflict/theme-monokai";
import 'ace-builds/src-noconflict/theme-mono_industrial';
import 'ace-builds/src-noconflict/theme-ambiance';
import 'ace-builds/src-noconflict/theme-clouds';
import 'ace-builds/src-noconflict/theme-cobalt';
import 'ace-builds/src-noconflict/theme-crimson_editor';
import 'ace-builds/src-noconflict/theme-dawn';
import 'ace-builds/src-noconflict/theme-dracula';
import 'ace-builds/src-noconflict/theme-dreamweaver';
import 'ace-builds/src-noconflict/theme-chaos';
import 'ace-builds/src-noconflict/theme-terminal';
import 'ace-builds/src-noconflict/theme-sqlserver';
import 'ace-builds/src-noconflict/theme-solarized_light';
import 'ace-builds/src-noconflict/theme-solarized_dark';
import 'ace-builds/src-noconflict/theme-pastel_on_dark';

config.setModuleUrl("ace/mode/java", modeJava);
config.setModuleUrl("ace/mode/c_cpp", modeCpp);
config.setModuleUrl("ace/mode/python", modePy);
config.setModuleUrl("ace/mode/rust", modeRust);
config.setModuleUrl("ace/mode/php", modePHP);
config.setModuleUrl("ace/mode/pascal", modePascal);

export default {
  components: {
    VAceEditor,
  },
  data() {
    return {
      content: JSON.stringify({ message: 'Hello Ace' }),
    };
  },
};
