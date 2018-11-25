import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from 'components/App';
import * as serviceWorker from './serviceWorker';
import {BrowserRouter} from 'react-router-dom';

import {client} from 'lib/apollo.js';
import {ApolloProvider} from 'react-apollo';

ReactDOM.render(
  <ApolloProvider client={client}>
    <BrowserRouter>
      <App />
    </BrowserRouter>
  </ApolloProvider>,
  document.getElementById('root')
);

// Register ServiceWorker for offline, fast reloads
serviceWorker.register();
