import axios from 'axios';
import { apiUrl } from '../axios';
import { language, token, userName } from '../account';

// 获取竞赛列表
export const getAllCompetitions = async () => {
    return await axios.get(`${apiUrl}/competitions`, {
        headers: {
            Username: encodeURIComponent(userName.value),
            Authorization: token.value,
        }
    })
}

// 获取指定竞赛ID信息
export const getCompetitionById = async (competitionId) => {
    return await axios.get(`${apiUrl}/competitions/info/${competitionId}`, {
        headers: {
            Username: encodeURIComponent(userName.value),
            Authorization: token.value,
        }
    })
}

// 加入竞赛
export const joinCompetition = async (competitionId) => {
    return await axios.post(`${apiUrl}/competitions/join/${competitionId}`, {}, {
        headers: {
            Username: encodeURIComponent(userName.value),
            Authorization: token.value,
            "Accept-Language": language.value
        }
    })
}

// 加入有密码的竞赛
export const joinCompWithPwd = async (competitionId, competitionPwd) => {
    return await axios.post(`${apiUrl}/competitions/join/pwd/${competitionId}`, {}, {
        params: {
            password: competitionPwd
        },
        headers: {
            Username: encodeURIComponent(userName.value),
            Authorization: token.value,
            "Accept-Language": language.value
        }
    })
}

// 退出竞赛
export const quitCompetition = async (competitionId) => {
    return await axios.post(`${apiUrl}/competitions/quit/${competitionId}`, {}, {
        headers: {
            Username: encodeURIComponent(userName.value),
            Authorization: token.value,
            "Accept-Language": language.value
        }
    })
}

// 查询用户是否在指定竞赛中
export const isInCompetition = async (competitionId) => {
    return await axios.get(`${apiUrl}/competitions/${competitionId}/in`, {
        headers: {
            Username: encodeURIComponent(userName.value),
            Authorization: token.value,
            "Accept-Language": language.value
        }
    })
}

// 获取指定竞赛的所有用户
export const getCompetitionUsers = async (competitionId) => {
    return await axios.get(`${apiUrl}/competitions/info/${competitionId}/users`, {
        headers: {
            Username: encodeURIComponent(userName.value),
            Authorization: token.value,
        }
    })
}

// 获取指定竞赛的所有题目
export const getCompetitionProblems = async (competitionId) => {
    return await axios.get(`${apiUrl}/competitions/info/${competitionId}/problems`, {
        headers: {
            Username: encodeURIComponent(userName.value),
            Authorization: token.value,
        }
    })
}