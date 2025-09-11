import axios from 'axios';
import { apiUrl, docsServer } from '../axios';
import { language, userName, token } from '../account';

// 获取公告
export const getAnnouncement = async () => {
    const resp = await fetch(`${docsServer}announcement.md`);
    return resp.text();
}

// 获取通知
export const getNotification = async () => {
    const resp = await fetch(`${docsServer}notice.md`);
    return resp.text();
}

// 更新用户简介
export const updateSynopsis = async (synopsis) => {
    const formData = new FormData();
    formData.append('synopsis', synopsis);
    return await axios.put(`${apiUrl}/users/synopsis`, formData, {
        headers: {
            Username: encodeURIComponent(userName.value),
            Authorization: token.value,
            "Accept-Language": language.value
        }
    })
}

// 修改头像
export const uploadAvatar = async (file) => {
    let formData = new FormData();
    formData.append('avatar', file);
    return await axios.put(`${apiUrl}/users/avatar`, formData, {
        headers: {
            'Content-Type': 'multipart/form-data',
            Username: encodeURIComponent(userName.value),
            Authorization: token.value,
            "Accept-Language": language.value
        },
    });
}

// 获取排行榜
export const getRanking = async () => {
    return await axios.get(`${apiUrl}/ranking`, {
        headers: {
            Username: encodeURIComponent(userName.value),
            Authorization: token.value,
            "Accept-Language": language.value
        }
    })
}