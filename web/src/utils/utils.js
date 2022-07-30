const normalizePath = path => path.replace(/[\\/]+/g, '/');

const buildPath = (...args) => {
    return args.map((part, i) => {
        if (i === 0) {
        return part.trim().replace(/[\/]*$/g, '')
        } else {
        return part.trim().replace(/(^[\/]*|[\/]*$)/g, '')
        }
    }).filter(x=>x.length).join('/')
}

export default { normalizePath, buildPath };