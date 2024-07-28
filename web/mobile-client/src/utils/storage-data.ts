interface storageType {
    key: string,
    value: {
        data: object,
        expired: number
    }
}

export const storageData = {
    //key保存键,value保存内容,*expired 失效时间,单位分钟
    set(key: string, value: string, expired: number): void {
        const now = Date.now();
        const source = { key: key, value: value };

        if (expired) {
            source.value = JSON.stringify({ data: value, expired: now + (1000 * 60 * expired) });
        } else {
            source.value = JSON.stringify({ data: value });
        }
        localStorage.setItem(source.key, source.value);
    },
    get(key: string): any {
        const now = Date.now();
        const source: storageType = {
            key: key,
            value: {
                data: {},
                expired: 0
            }
        };

        const readStorage = localStorage.getItem(source.key);
        source.value = JSON.parse(readStorage ? readStorage : "{}");

        if (source.value) {
            //超过失效时 删除
            if (!source.value.expired) {
                return source.value.data;
            }
            else if (now >= source.value.expired) {
                this.remove(source.key);
                return null;
            } else {
                return source.value.data;
            }
        }

        return null;
    },
    //更新
    update(key: string, data: any) {
        let read = localStorage.getItem(key);
        if (read) {
            var value = JSON.parse(read);
            value.data = data;
            localStorage.setItem(key, JSON.stringify(value));
        }
    },
    remove(key: string): void {
        localStorage.removeItem(key);
    },
}