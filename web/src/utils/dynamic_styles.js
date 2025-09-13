// 根据结果不同显示不同颜色
export const getResultStyle = (result) => {
    switch (result) {
        case 'Compile Failed':
            return 'color: red; font-weight: bold;';
        case 'Time Limit Exceeded':
            return 'color: red; font-weight: bold;';
        case 'Memory Limit Exceeded':
            return 'color: red; font-weight: bold;';
        case 'Internal Error':
            return 'color: red; font-weight: bold';
        case 'Accepted':
            return 'color: green; font-weight: bold;';
        case 'Failed':
            return 'color: orange; font-weight: bold;';
        case 'Wrong Answer':
            return 'color: orange; font-weight: bold;';
        default:
            return '';
    }
};

// 根据结果不同显示不同Chip颜色
export const getResultChipColor = (result) => {
  switch (result) {
    case 'Accepted':
      return 'success';
    case 'Wrong Answer':
      return 'error';
    case 'Time Limit Exceeded':
      return 'warning';
    case 'Memory Limit Exceeded':
      return 'warning';
    case 'Compile Failed':
      return 'error';
    case 'Internal Error':
      return 'warning';
    case 'Failed':
      return 'error';
    default:
      return 'default';
  }
}

// 根据题目难度显示不同字体
export const difficultyColor = (difficulty) => {
    switch (difficulty) {
        case 0:
            return 'font-weight: bold;color: green;';
        case 1:
            return 'font-weight: bold;color: orange;';
        case 2:
            return 'font-weight: bold;color: red;';
        default:
            return 'font-weight: bold;color: green;';
    }
};

export const difficultyLang = (difficulty) => {
    switch (difficulty) {
        case 0:
            return 'message.easy';
        case 1:
            return 'message.medium';
        case 2:
            return 'message.hard';
        default:
            return 'message.easy';
    }
}