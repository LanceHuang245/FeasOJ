import axios from 'axios';
import { apiUrl } from '../axios';
import { language, userName, token } from '../account';

// 获取讨论列表
export const getAllDis = async (page, itemsPerPage) => {
    return await axios.get(`${apiUrl}/discussions`, {
        params: {
            page: page,
            itemsPerPage: itemsPerPage
        },
        headers: {
            Username: encodeURIComponent(userName.value),
            Authorization: token.value
        }
    })
}

// 获取讨论详细信息
export const getDisDetails = async (did) => {
    return await axios.get(`${apiUrl}/discussions/${did}`, {
        headers: {
            Username: encodeURIComponent(userName.value),
            Authorization: token.value
        }
    })
}

// 获取指定讨论的所有回复
export const getComments = async (did) => {
    return await axios.get(`${apiUrl}/discussions/comments/${did}`, {
        headers: {
            Username: encodeURIComponent(userName.value),
            Authorization: token.value
        }
    })
}

// 添加讨论
export const addDiscussion = async (Title, Content) => {
    const formData = new FormData();
    formData.append('title', Title);
    formData.append('content', Content);
    return await axios.post(`${apiUrl}/discussions/add`, formData, {
        headers: {
            Username: encodeURIComponent(userName.value),
            Authorization: token.value,
            "Accept-Language": language.value
        },
    });
}

// 添加评论
export const addComment = async (did, content) => {
    const formData = new FormData();
    formData.append('content', content);
    return await axios.post(`${apiUrl}/discussions/comments/add/${did}`, formData, {
        headers: {
            Username: encodeURIComponent(userName.value),
            Authorization: token.value,
            "Accept-Language": language.value
        },
    })
}

// 删除讨论
export const deleteDiscussion = async (id) => {
    return await axios.post(`${apiUrl}/discussions/delete/${id}`, {}, {
        headers: {
            Username: encodeURIComponent(userName.value),
            Authorization: token.value,
            "Accept-Language": language.value
        }
    })
}

// 删除讨论评论
export const deleteComment = async (id) => {
    return await axios.post(`${apiUrl}/discussions/comments/delete/${id}`, {}, {
        headers: {
            Username: encodeURIComponent(userName.value),
            Authorization: token.value,
            "Accept-Language": language.value
        },
    })
}