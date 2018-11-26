import React from 'react';
import styles from './App.module.css';
import logoUrl from 'assets/beatman.svg';
import {NavLink, Switch, Route} from 'react-router-dom';
import Song from 'components/Song.js';
import gql from 'graphql-tag';
import {Query} from 'react-apollo';

const App = () => (
  <div className={styles.app}>
    <div className={styles.header}>
      <img src={logoUrl} alt="BeatMan" className={styles.logo} />
      <NavLink className={styles.nav} activeClassName={styles.navActive} exact to="/">
        Library
      </NavLink>
      <NavLink className={styles.nav} activeClassName={styles.navActive} to="/downloads">
        Downloads
      </NavLink>
    </div>

    <div className={styles.content}>
      <Switch>
        <Route path="/downloads">{() => <h1>Downloads</h1>}</Route>

        {/* TODO: replace this with an actual implementation */}
        <Route exact path="/">
          {() => (
            <div>
              <h1>Library</h1>
              <Query
                query={gql`
                  query {
                    songs {
                      id
                      ...Song
                    }
                  }
                  ${Song.fragment}
                `}
              >
                {({loading, error, data}) => {
                  if (loading) {
                    return <div>Loading...</div>;
                  }
                  return data.songs.map(e => <Song key={e.id} {...e} />);
                }}
              </Query>
            </div>
          )}
        </Route>
      </Switch>
    </div>
  </div>
);

export default App;
