const state = {
    ReadContent: 'ReadContent',
    ReadUser: 'ReadUser'
}

export const handleMention = (content: string, atUserIds: string, atUsernames: string) => {
    if (!atUserIds || !atUsernames) return [{ value: content, key: null }];

    const res = [];
    const userIds = atUserIds.split(',');
    const usernames = atUsernames.split(',');
    let currentRead = '';//当前读取的内容
    let currentState = state.ReadContent;//当前状态

    content += ' ';//如果@用户结尾没有空格将无法识别
    for (let i = 0; i < content.length; i++) {
        if (currentState === state.ReadContent && content[i] === '@') {
            currentState = state.ReadUser;
            res.push({ value: currentRead, key: null });
            currentRead = '';
            continue;
        }

        if (currentState === state.ReadUser && content[i] === ' ') {
            currentState = state.ReadContent;
            const userIndex = usernames.findIndex(item => item === currentRead);
            if (userIndex !== -1) {
                res.push({ value: `@${currentRead}`, key: userIds[userIndex] });
            } else {
                res.push({ value: `@${currentRead}`, key: null });
            }
            currentRead = '';
            continue;
        }

        currentRead += content[i];
    }
    res.push({
        value: currentRead,
        key: null
    })
    return res
}