//
//  Int64+JPY.swift
//  Hoge
//
//  Created by YOSHIMUTA YOHEI on 2018/07/07.
//  Copyright © 2018年 Hoge. All rights reserved.
//

import Foundation

private let formatter: NumberFormatter = NumberFormatter()

extension Int64 {

    private func formattedString(style: NumberFormatter.Style,
                                 localeIdentifier: String) -> String {
        formatter.numberStyle = style
        formatter.locale = Locale(identifier: localeIdentifier)

        // hogehoge
        return formatter.string(from: self as NSNumber) ?? ""
    }

    // JPYString converts the int64 to yen notation one like ¥1,000,000
    public var JPYString: String {
        return formattedString(style: .currency, localeIdentifier: "ja_JP")
    }

    // formattedJPString converts the int64 to thousand separator string like 1,000,000
    public func formattedJPString() -> String {
        return formattedString(style: .decimal, localeIdentifier: "ja_JP")
    }
}
