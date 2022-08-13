let slogan = [
        // ['截止时间', '口号']
        // ['2022-06-30 23:59:59', '6月最后的战斗，我要全力以赴，使命必达！'],
        // ['2022-07-31 23:59:59', '7月，为自己而战，全力以赴每一天！'],
        ['2022-08-20 23:59:59', '坚持拜访量达标，你才知道什么是惊喜！~'],
        ['2022-08-31 23:59:59', '月底冲刺，你的团队为你骄傲，酵母为你骄傲。'],
        ['2022-12-31 23:59:59', '浪是一种气质，量是一种态度，2022下半年浪起来啦！']
];

function getSlogan(ts) {
        let tsInMs = ts * 1000
        for (let i in slogan) {
                let item = slogan[i]
                let endTime = Date.parse(item[0])
                if (tsInMs <= endTime) {
                        return item[1];
                }
        }
        return ''
}

