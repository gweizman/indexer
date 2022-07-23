import React, {useState} from 'react';
import {Treebeard, decorators} from 'react-treebeard';
import Header from './Header';

const ddata = {
    name: '.',
    id: ".",
    toggled: false,
    loading: true,
    children: []
};

const TreeExample = () => {
    const [data, setData] = useState(ddata);
    const [cursor, setCursor] = useState(false);
    
    const onToggle = (node, toggled) => {
        if (cursor) {
            cursor.active = false;
        }
        node.active = true;
        if (node.children) {
            node.toggled = toggled;
        }
        console.log(node.id)
        setCursor(node);
        setData(Object.assign({}, data))
    }
    
    return (
        <Treebeard 
            data={data}
            onToggle={onToggle}
            decorators={{...decorators, Header}}
            customStyles={{
                header: {
                    title: {
                        color: 'red'
                    }
                }
            }}
        />
    )
}

function FileTree() {
    return (
        <TreeExample/>
    )
}

export default FileTree;