import { useQuery } from '@tanstack/react-query'
import React, {useState} from 'react';
import { useNavigate } from "react-router-dom";
import utils from '../utils/utils';
import "./FileTree.css"

function isSubPath(path1, path2, exact) {
    let tree1 = path1.split(/[\\/]/)
    let tree2 = path2.split(/[\\/]/)

    if (tree1.length > tree2.length) {
        return false;
    }
    if (exact && tree1.length !== tree2.length) {
        return false;
    }

    for (let i = 0; i < tree1.length; i++) {
        if (tree1[i] != tree2[i]) {
            return false;
        }
    }

    return true;
}

function TreeDir(props) {
    const [toggleValue, setToggleValue]  = useState(isSubPath(props.path, props.current, false))

    const handleClick = (event) => {
        setToggleValue(!toggleValue);
        event.preventDefault();
    }

    return (
        <div>
            <div className="tree-item" onClick={handleClick}>
                {toggleValue ? <i className="fa fa-chevron-down"></i> : <i className="fa fa-chevron-right"></i>}
                {" "}
                <i className="fa fa-folder"></i>
                {" "}
                {props.path.split(/[\\/]/).slice(-1)[0]}
            </div>
            <div className="tree-children">{ toggleValue ? <TreeChildren project={props.project} parent={props.path} current={props.current} /> : <div></div>}</div>
        </div>
    )
}

function TreeFile(props) {
    const navigate = useNavigate()
    const handleClick = (e) => {
        navigate({
            pathname:utils.buildPath("/browse", props.project, props.path, props.name),
          });
          e.preventDefault();
    }

    const selected = isSubPath(utils.buildPath(props.path, props.name), props.current, true);

    return (
        <div className={"tree-item " + (selected ? "selected" : "")} onClick={handleClick}>
            <i className="fa fa-file-text"></i>
            {" "}
            {props.name}
        </div>
    )
}

function TreeChildren(props) {
    let { isLoading, error, data } = useQuery(['dirList', props.parent], () =>
        fetch(utils.buildPath('/api/OpenGrok/dir/', props.parent)).then(res =>
            res.json()
        )
    )
    if (isLoading) return 'Loading...'
    if (error) return 'An error has occurred: ' + error.message
    console.log(props.parent)

    data = data.sort((a, b) => {
        if (a.is_dir && !b.is_dir) {
            return -1
        }
        if (b.is_dir && !a.is_dir) {
            return 1
        }

        if (a.is_dir) {
            return a.path.localeCompare(b.path)
        }
        return a.name.localeCompare(b.name)
    })

    return (
        <div>
            {
                data.map((object, i) => {
                    if (object.is_dir) {
                        return <TreeDir project={props.project} path={object.path} current={props.current}/>
                    }
                    return <TreeFile project={props.project} path={object.path} name={object.name} current={props.current} />
                })
            }
        </div>
    )
}

function FileTree(props) {
    return (
        <div style={{height: "100%"}}>
            <strong>{props.project}</strong>
            <div className="file-tree">
                <TreeChildren parent="." project={props.project} current={props.current} />
            </div>
        </div>
    )
}

export default FileTree;