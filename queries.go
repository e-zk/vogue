package vogue

// graphql queries
const (
	qBrands      = "query{allBrands{Brand{name slug}}}"
	qSeasons     = "query{allSeasons{Season{name slug}}}"
	qSeasonShows = `query{ allContent( type: ["FashionShowV2"], first: 1000, filter: { season: { slug: "%s" } }) { Content { id GMTPubDate url title slug ... on FashionShowV2 { instantShow brand { name slug } season { name slug year } photosTout { ... on Image { url } } } } } }`
	qBrandShows  = `query { allContent(type: ["FashionShowV2"], first: 1000, filter: { brand: { slug: "%s" } }) { Content { id GMTPubDate url title slug ... on FashionShowV2 { instantShow brand { name slug } season { name slug year } photosTout { ... on Image { url } } } } } }`
	qFashionShow = `query {
  fashionShowV2(slug: "%s") {
    GMTPubDate
    url
    title
    slug
    id
    city {
      name
    }
    brand {
      name
      slug
    }
    season {
      name
      slug
      year
    }
    photosTout {
      ... on Image {
        url
      }
    }
    review {
      pubDate
      body
      contributor {
        author {
          name
        }
      }
    }
    galleries {
      collection {
        ...GalleryFragment
      }
      atmosphere {
        ...GalleryFragment
      }
      beauty {
        ...GalleryFragment
      }
      detail {
        ...GalleryFragment
      }
      frontRow {
        ...GalleryFragment
      }
    }
    video {
      url
      cneId
      title
    }
  }
}
fragment GalleryFragment on FashionShowGallery {
  title
  slidesV2 {
    ... on GallerySlidesConnection {
      slide {
        ... on Slide {
          id
          credit
          photosTout {
            ...imageFields
          }
        }
        ... on CollectionSlide {
          id
          type
          credit
          title
          photosTout {
            ...imageFields
          }
        }
        __typename
      }
    }
  }
}
fragment imageFields on Image {
  id
  url
  caption
  credit
  width
  height
}
`
)
