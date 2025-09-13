import axios from 'axios';
import { apiUrl } from '../axios';
import { language, userName, token } from '../account';

// 管理员获取竞赛列表
export const getAllCompetitionsInfo = async () => {
    return await axios.get(`${apiUrl}/admin/competitions`, {
        headers: {
            Username: encodeURIComponent(userName.value),
            Authorization: token.value,
            "Accept-Language": language.value
        }
    })
}

// 管理员获取指定列表信息
export const getCompetitionInfoByIDAdmin = async (cid) => {
    return await axios.get(`${apiUrl}/admin/competitions/${cid}`, {
        headers: {
            Username: encodeURIComponent(userName.value),
            Authorization: token.value
        }
    })
}

// 管理员获取指定题目所有信息
export const getProblemAllInfoByAdmin = async (pid) => {
    return await axios.get(`${apiUrl}/admin/problems/${pid}`, {
        headers: {
            Username: encodeURIComponent(userName.value),
            Authorization: token.value
        },
    })
}

// 管理员获取所有用户信息
export const getAllUsersInfo = async () => {
    return await axios.get(`${apiUrl}/admin/users`, {
        headers: {
            Username: encodeURIComponent(userName.value),
            Authorization: token.value
        }
    })
}

// 管理员添加/更新题目信息
export const updateProblemInfo = async (problemInfo) => {
    return await axios.post(`${apiUrl}/admin/problems`, problemInfo, {
        headers: {
            Username: encodeURIComponent(userName.value),
            Authorization: token.value,
            "Content-Type": "application/json",
            "Accept-Language": language.value
        }
    })
}

// 管理员删除题目信息
export const deleteProblemAllInfo = async (pid) => {
    return await axios.delete(`${apiUrl}/admin/problems/${pid}`, {
        headers: {
            Username: encodeURIComponent(userName.value),
            Authorization: token.value,
            "Accept-Language": language.value
        }
    })
}

// 封禁用户
export const banUser = async (uid) => {
    return await axios.put(`${apiUrl}/admin/users/status`, {}, {
        params: {
            user_id: uid,
            status: true
        },
        headers: {
            Username: encodeURIComponent(userName.value),
            Authorization: token.value,
            "Accept-Language": language.value
        }
    })
}

// 解封用户
export const unbanUser = async (uid) => {
    return await axios.put(`${apiUrl}/admin/users/status`, {}, {
        params: {
            user_id: uid,
            status: false
        },
        headers: {
            Username: encodeURIComponent(userName.value),
            Authorization: token.value,
            "Accept-Language": language.value
        }
    })
}

// 晋升用户
export const promoteUser = async (uid) => {
    return await axios.put(`${apiUrl}/admin/users/promote`, {}, {
        params: {
            uid: uid
        },
        headers: {
            Username: encodeURIComponent(userName.value),
            Authorization: token.value,
            "Accept-Language": language.value
        }
    })
}

// 降级用户
export const demoteUser = async (uid) => {
    return await axios.put(`${apiUrl}/admin/users/demote`, {}, {
        params: {
            uid: uid
        },
        headers: {
            Username: encodeURIComponent(userName.value),
            Authorization: token.value,
            "Accept-Language": language.value
        }
    })
}

// 管理员获取题目列表
export const getAllProblemsAdmin = async () => {
    return await axios.get(`${apiUrl}/admin/problems`, {
        headers: {
            Username: encodeURIComponent(userName.value),
            Authorization: token.value,
        }
    });
}

// 管理员删除竞赛
export const deleteCompetition = async (cid) => {
    return await axios.delete(`${apiUrl}/admin/competitions/${cid}`, {
        headers: {
            Username: encodeURIComponent(userName.value),
            Authorization: token.value,
            "Accept-Language": language.value
        }
    })
}

// 管理员添加/更新竞赛信息
export const updateComInfo = async (comInfo) => {
    return await axios.post(`${apiUrl}/admin/competitions`, comInfo, {
        headers: {
            Username: encodeURIComponent(userName.value),
            Authorization: token.value,
            "Content-Type": "application/json",
            "Accept-Language": language.value
        }
    })
}

// 管理员启用竞赛计分
export const caculateComScore = async (cid) => {
    return await axios.get(`${apiUrl}/admin/competitions/${cid}/score`, {
        headers: {
            Username: encodeURIComponent(userName.value),
            Authorization: token.value,
            "Accept-Language": language.value
        }
    })
}

// 管理员查看竞赛得分情况
export const getScores = async (cid, page, itemsPerPage) => {
    return await axios.get(`${apiUrl}/admin/competitions/${cid}/scoreboard`, {
        params: {
            page: page,
            itemsPerPage: itemsPerPage
        },
        headers: {
            Username: encodeURIComponent(userName.value),
            Authorization: token.value,
            "Accept-Language": language.value
        }
    })
}

// 管理员获取IP统计
export const getIpStat = async () => {
    return await axios.get(`${apiUrl}/admin/ipstats`, {
        headers: {
            Username: encodeURIComponent(userName.value),
            Authorization: token.value,
        }
    })
}