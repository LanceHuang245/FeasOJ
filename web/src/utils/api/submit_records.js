import axios from 'axios';
import { apiUrl } from '../axios';
import { userName, token } from '../account';

// 获取指定用户提交记录
export const getUserSubmitRecords = async (username) => {
    return await axios.get(`${apiUrl}/users/${username}/submitrecords`,{
        headers: {
            Username: encodeURIComponent(userName.value),
            Authorization: token.value
        }
    })
}

// 获取30天内的提交记录
export const getSubmitRecords = async () => {
    return await axios.get(`${apiUrl}/submitrecords`, {
        headers: {
            Username: encodeURIComponent(userName.value),
            Authorization: token.value
        }
    })
}