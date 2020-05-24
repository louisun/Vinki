import React from 'react';

class Latex extends React.Component {
    constructor(props) {
        super(props);
        this.node = React.createRef();
    }

    componentDidMount() {
        this.renderMath();
    }

    componentDidUpdate() {
        this.renderMath();
    }

    renderMath() {
        window.MathJax.Hub.Queue([
            "Typeset",
            window.MathJax.Hub,
            this.node.current
        ]);
    }

    render() {
        // const { text } = this.props;
        return <div ref={this.node}>{this.props.children}</div>;
    }
}

export default Latex;