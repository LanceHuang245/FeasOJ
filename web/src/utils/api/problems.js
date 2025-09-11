import axios from 'axios';
import { apiUrl } from '../axios';
import { language, userName, token } from '../account';

// 获取题目列表
export const getAllProblems = async () => {
    return await axios.get(`${apiUrl}/problems`, {
        headers: {
            Username: encodeURIComponent(userName.value),
            Authorization: token.value
        }
    });
}

// 获取题目详细信息
export const getPbDetails = async (pid) => {
    return await axios.get(`${apiUrl}/problems/${pid}`, {
        headers: {
            Username: encodeURIComponent(userName.value),
            Authorization: token.value
        }
    });
}

// 提交代码文件
export const uploadCode = async (file, pid) => {
    let formData = new FormData();
    formData.append('code', file);
    return await axios.post(`${apiUrl}/problems/${pid}/code`, formData, {
        headers: {
            'Content-Type': 'multipart/form-data',
            Username: encodeURIComponent(userName.value),
            Authorization: token.value,
            "Accept-Language": language.value
        },
    });
}